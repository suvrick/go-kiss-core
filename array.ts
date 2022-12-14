import bigInt from "big-integer";

class ArrayBufferPool
{
	private static buffers: Uint8ClampedArray[][] = [];

	public static add(buffer: Uint8ClampedArray): void
	{
		if (!(buffer.byteLength in ArrayBufferPool.buffers))
			ArrayBufferPool.buffers[buffer.byteLength] = [];

		ArrayBufferPool.buffers[buffer.byteLength].push(buffer);
	}

	public static get(size: number): Uint8ClampedArray
	{
		if (!(size in ArrayBufferPool.buffers))
			ArrayBufferPool.buffers[size] = [];

		if (ArrayBufferPool.buffers[size].length == 0)
			return new Uint8ClampedArray(new ArrayBuffer(size));

		const len: number = ArrayBufferPool.buffers[size].length - 1;
		const remove: Uint8ClampedArray = ArrayBufferPool.buffers[size][len];
		ArrayBufferPool.buffers[size].splice(len, 1);
		return remove;
	}
}

export class ByteArray
{
	private static readonly STEP: number = 65000;
	private static readonly SIZE_OF_INT8: number = 1;

	private bufferU8: Uint8ClampedArray;
	private dataView: DataView;

	public length: number = 0;
	public position: number = 0;

	public constructor(buffer: ArrayBuffer | null = null)
	{
		if (buffer == null)
			this.bufferU8 = ArrayBufferPool.get(1024);
		else
		{
			this.bufferU8 = new Uint8ClampedArray(buffer);
			this.length = buffer.byteLength;
		}

		this.dataView = new DataView(this.bufferU8.buffer);
	}

	public get bytesAvailable(): number
	{
		return this.length - this.position;
	}

	public setArrayBuffer(newBuffer: ArrayBuffer): void
	{
		this.bufferU8 = new Uint8ClampedArray(newBuffer);
		this.position = 0;
		this.length = newBuffer.byteLength;

		this.dataView = new DataView(this.bufferU8.buffer);
	}

	public readBytes(bytes: ByteArray, length: number): void
	{
		const savePos: number = bytes.position;

		bytes.writeBytes(this, this.position, length);

		this.position += length;
		bytes.position = savePos;
	}

	public writeIntLeb(value: number): number
	{
		if (value < 0)
			value = value >>> 0;

		const encodedNum: number[] = ByteArray.encodeLeb128(value);

		for (const part of encodedNum)
			this.writeByte(part);

		return encodedNum.length;
	}

	public writeBytes(bytes: ByteArray, offset: number = 0, length: number = 0): void
	{
		if (length == 0)
			length = bytes.length - offset;

		let byteLength: number = this.bufferU8.byteLength >>> 0;
		if (this.length < this.position + length)
			this.length = this.position + length;

		while (byteLength < this.length)
			byteLength *= 2;

		if (byteLength == this.bufferU8.byteLength)
			this.bufferU8.set(new Uint8ClampedArray(bytes.bufferU8.buffer, offset, length), this.position);
		else
		{
			const buffer: Uint8ClampedArray = ArrayBufferPool.get(byteLength);
			buffer.set(this.bufferU8);
			buffer.set(new Uint8ClampedArray(bytes.bufferU8.buffer, offset, length), this.position);
			ArrayBufferPool.add(this.bufferU8);
			this.bufferU8 = buffer;

			this.dataView = new DataView(this.bufferU8.buffer);
		}

		this.position += length;
	}

	public readByte(): number
	{
		if (this.position + ByteArray.SIZE_OF_INT8 > this.length)
			throw new Error("Failed to read past end of the stream");

		const r: number = this.dataView.getInt8(this.position);
		this.position++;
		return r;
	}

	public readUnsignedByte(): number
	{
		if (this.position + ByteArray.SIZE_OF_INT8 > this.length)
			throw new Error("Failed to read past end of the stream");

		const r: number = this.dataView.getUint8(this.position);
		this.position++;
		return r;
	}

	public readIntLeb(): number
	{
		let result: bigInt.BigInteger = bigInt(0);

		let value: bigInt.BigInteger = bigInt(0);
		let shift: bigInt.BigInteger = bigInt(0);

		do
		{
			value = bigInt(this.dataView.getUint8(this.position));
			this.position++;
			let _byte: bigInt.BigInteger = value.and(bigInt(0x7F));

			_byte = _byte.shiftLeft(bigInt(shift));

			result = result.or(_byte);
			shift = shift.plus(bigInt(7));
		}
		while (value.greaterOrEquals(bigInt(128)));

		if (result.greaterOrEquals(Number.MAX_SAFE_INTEGER))
			return eval("String(result)");

		return result.toJSNumber();
	}

	public readLebUTF(): string
	{
		return this.readUTFBytes(this.readIntLeb());
	}

	public readFloat(): number
	{
		const r: number = this.dataView.getFloat32(this.position, true);
		this.position += 4;
		return r;
	}

	public toString(): string
	{
		const str: string = this.fromCharCode();
		try
		{
			return decodeURIComponent(escape(str));
		}
		catch (e)
		{
			return str;
		}
	}

	public clear(): void
	{
		this.length = 0;
		this.position = 0;
	}

	public get(key: number): number
	{
		return this.bufferU8[key];
	}

	public set(key: number, value: number): void
	{
		if (key == this.length)
		{
			this.length++;
			this.position++;
		}

		while (this.length >= this.bufferU8.byteLength)
			this.expand();

		this.bufferU8[key] = value;
	}

	protected writeUTFLeb(value: string): void
	{
		try
		{
			value = unescape(encodeURIComponent(value));
			this.writeIntLeb(value.length);
			this.writeUTFBytes(value);
		}
		catch (e)
		{
			if (e instanceof URIError)
				console.error("URIError with value:", value);
		}
	}

	protected writeByte(value: number): void
	{
		if (this.length < this.position + ByteArray.SIZE_OF_INT8)
			this.length = this.position + ByteArray.SIZE_OF_INT8;

		while (this.length >= (this.bufferU8.byteLength >>> 0))
			this.expand();

		this.dataView.setUint8(this.position, (value & 0xff));
		this.position++;
	}

	protected write(length: number): void
	{
		if (this.length < this.position + (length >>> 0))
			this.length = this.position + (length >>> 0);

		while (this.length >= this.bufferU8.byteLength)
			this.expand();
	}

	private writeUTFBytes(value: string): void
	{
		let i: number = 0;
		const l: number = value.length;

		for (; i < l - 10; i += 10)
		{
			this.writeByte(value.charCodeAt(i));
			this.writeByte(value.charCodeAt(i + 1));
			this.writeByte(value.charCodeAt(i + 2));
			this.writeByte(value.charCodeAt(i + 3));
			this.writeByte(value.charCodeAt(i + 4));
			this.writeByte(value.charCodeAt(i + 5));
			this.writeByte(value.charCodeAt(i + 6));
			this.writeByte(value.charCodeAt(i + 7));
			this.writeByte(value.charCodeAt(i + 8));
			this.writeByte(value.charCodeAt(i + 9));
		}

		for (; i < l; i++)
			this.writeByte(value.charCodeAt(i));
	}

	private readUTFBytes(length: number): string
	{
		const step: number = ByteArray.STEP;
		const _length: number = (length + this.position);

		let str: string | null = null;
		for (let i = this.position; i < _length; i += step)
		{
			const size: number = Math.min(step, _length - i);

			// @ts-ignore
			// eslint-disable-next-line prefer-spread
			const s: string = String.fromCharCode.apply(String, new Uint8ClampedArray(this.bufferU8.buffer, i, size));
			if (str == null)
				str = s;
			else
				str += s;
		}
		if (str == null)
			str = "";

		const result: string = decodeURIComponent(escape(str));

		this.position += length;
		return result;
	}

	private expand(): void
	{
		const newBuffer: Uint8ClampedArray = ArrayBufferPool.get(this.bufferU8.byteLength * 2);
		newBuffer.set(this.bufferU8);

		ArrayBufferPool.add(this.bufferU8);
		this.bufferU8 = newBuffer;

		this.dataView = new DataView(this.bufferU8.buffer);
	}

	private fromCharCode(): string
	{
		let str: string = "";
		const step: number = ByteArray.STEP;
		for (let i = 0; i < this.length; i += step)
		{
			const size: number = Math.min(step, this.length - i);

			// @ts-ignore
			// eslint-disable-next-line prefer-spread
			str += String.fromCharCode.apply(String, new Uint8ClampedArray(this.bufferU8.buffer, i, size));
		}
		return str;
	}

	public static encodeLeb128(value: number): number[]
	{
		const list: number[] = [];

		let newValue: bigInt.BigInteger = bigInt(value);

		do
		{
			let _byte: bigInt.BigInteger = newValue.and(bigInt(0x7f));

			newValue = newValue.shiftRight(bigInt(7));

			if (newValue.notEquals(bigInt(0)))
				_byte = _byte.or(bigInt(0x80));
			list.push(Number(_byte));
		}
		while (newValue.notEquals(0));

		return list;
	}
}