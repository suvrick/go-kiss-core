import { ByteArray } from "../framework/net/ByteArray";

export class PacketParser
{
	private static split: boolean = true;

	public static readData(buffer: ByteArray, format: string, output: any[], type: number): void
	{
		let optional: boolean = false;

		const groups: any[] = [ output ];

		for (let i: number = 0; i < format.length; i++)
		{
			const symbol: string = format.charAt(i);
			const groupCur: any[] = groups[groups.length - 1];

			if (symbol === ",")
			{
				if (optional || groupCur !== output)
					throw new Error("Bad signature 2 for server packet");

				optional = true;
				continue;
			}

			if (symbol === "]")
			{
				if (groupCur === output)
					throw new Error("Bad signature 3 for server packet");

				groupCur["group_length"]--;

				if (groupCur["group_length"] !== 0)
				{
					i = groupCur["group_pos"];
					continue;
				}
				PacketParser.splitGroup(groupCur);
				delete groupCur["group_length"];
				delete groupCur["group_pos"];
				delete groupCur["group_total"];

				groups.pop();
				continue;
			}

			if (buffer.bytesAvailable === 0)
			{
				if (optional && groupCur === output)
					break;
				throw new Error("No data for server packet " + type);
			}

			if (symbol === "[")
			{
				const groupNew: any[] = [];
				groupNew["group_total"] = groupNew["group_length"] = buffer.readIntLeb();

				groupCur.push(groupNew);

				if (groupNew["group_length"] !== 0)
				{
					groupNew["group_pos"] = i;
					groups.push(groupNew);
					continue;
				}

				i = PacketParser.getGroupLast(format, i);
				continue;
			}
			switch (symbol)
			{
				case "A":
				{
					const length: number = buffer.readIntLeb();
					const array: ByteArray = new ByteArray();
					if (length !== 0)
						buffer.readBytes(array, 0, length);
					groupCur.push(array);
					break;
				}
				case "S":
					groupCur.push(buffer.readLebUTF());
					break;
				case "F":
					groupCur.push(buffer.readFloat());
					break;
				case "I":
					groupCur.push(buffer.readIntLeb());
					break;
				case "B":
					groupCur.push(buffer.readUnsignedByte());
					break;
			}
		}

		if (buffer.bytesAvailable === 0)
			return;

		throw new Error("Data left in server packet " + type + " Avalible: " + buffer.bytesAvailable);
	}

	private static splitGroup(group: object[]): void
	{
		if (group["group_total"] === group.length || !this.split)
			return;

		const itemLength: number = group.length / group["group_total"];
		const result: object[] = [];
		const total: number = group["group_total"];

		for (let i = 0; i < total; i++)
			result.push(group.slice(i * itemLength, (i + 1) * itemLength));

		group.length = 0;
		group.push(...result);
	}

	private static getGroupLast(format: string, first: number): number
	{
		let left: number = 1;
		for (let last: number = first + 1; last < format.length; last++)
		{
			switch (format.charAt(last))
			{
				case "]":
					left--;
					break;
				case "[":
					left++;
					break;
			}

			if (left !== 0)
				continue;
			if (last === first + 1)
				break;
			return last;
		}

		throw new Error("Bad signature 1 for server packet");
	}
}