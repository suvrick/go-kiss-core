!(function (t) {
    var e = {};
    function r(n) {
        if (e[n]) return e[n].exports;
        var o = (e[n] = { i: n, l: !1, exports: {} });
        return t[n].call(o.exports, o, o.exports, r), (o.l = !0), o.exports;
    }
    (r.m = t),
        (r.c = e),
        (r.d = function (t, e, n) {
            r.o(t, e) || Object.defineProperty(t, e, { enumerable: !0, get: n });
        }),
        (r.r = function (t) {
            "undefined" != typeof Symbol && Symbol.toStringTag && Object.defineProperty(t, Symbol.toStringTag, { value: "Module" }), Object.defineProperty(t, "__esModule", { value: !0 });
        }),
        (r.t = function (t, e) {
            if ((1 & e && (t = r(t)), 8 & e)) return t;
            if (4 & e && "object" == typeof t && t && t.__esModule) return t;
            var n = Object.create(null);
            if ((r.r(n), Object.defineProperty(n, "default", { enumerable: !0, value: t }), 2 & e && "string" != typeof t))
                for (var o in t)
                    r.d(
                        n,
                        o,
                        function (e) {
                            return t[e];
                        }.bind(null, o)
                    );
            return n;
        }),
        (r.n = function (t) {
            var e =
                t && t.__esModule
                    ? function () {
                          return t.default;
                      }
                    : function () {
                          return t;
                      };
            return r.d(e, "a", e), e;
        }),
        (r.o = function (t, e) {
            return Object.prototype.hasOwnProperty.call(t, e);
        }),
        (r.p = ""),
        r((r.s = 2));
})([
    function (module, exports, __webpack_require__) {
        "use strict";
        var __importDefault =
            (this && this.__importDefault) ||
            function (t) {
                return t && t.__esModule ? t : { default: t };
            };
        Object.defineProperty(exports, "__esModule", { value: !0 }), (exports.ByteArray = void 0);
        const big_integer_1 = __importDefault(__webpack_require__(12));
        class ArrayBufferPool {
            static add(t) {
                t.byteLength in ArrayBufferPool.buffers || (ArrayBufferPool.buffers[t.byteLength] = []), ArrayBufferPool.buffers[t.byteLength].push(t);
            }
            static get(t) {
                if ((t in ArrayBufferPool.buffers || (ArrayBufferPool.buffers[t] = []), 0 === ArrayBufferPool.buffers[t].length)) return new Uint8ClampedArray(new ArrayBuffer(t));
                const e = ArrayBufferPool.buffers[t].length - 1,
                    r = ArrayBufferPool.buffers[t][e];
                return ArrayBufferPool.buffers[t].splice(e, 1), r;
            }
        }
        ArrayBufferPool.buffers = [];
        class ByteArray {
            constructor(t = null) {
                (this._length = 0),
                    (this._position = 0),
                    null == t ? (this.bufferU8 = ArrayBufferPool.get(1024)) : ((this.bufferU8 = new Uint8ClampedArray(t)), (this.length = t.byteLength)),
                    (this.dataView = new DataView(this.bufferU8.buffer));
            }
            get length() {
                return this._length;
            }
            set length(t) {
                this._length = t;
            }
            get bytesAvailable() {
                return this.length - this.position;
            }
            get position() {
                return this._position;
            }
            set position(t) {
                this._position = t;
            }
            setArrayBuffer(t) {
                (this.bufferU8 = new Uint8ClampedArray(t)), (this.position = 0), (this.length = t.byteLength), (this.dataView = new DataView(this.bufferU8.buffer));
            }
            readBytes(t, e = 0, r = 0) {
                const n = t.position;
                (t.position = e), 0 === r && (r = this._length - e), t.writeBytes(this, this.position, r), (this.position += r), (t.position = n);
            }
            writeIntLeb(t) {
                t < 0 && (t >>>= 0);
                const e = ByteArray.encodeLeb128(t);
                for (const t of e) this.writeByte(t);
                return e.length;
            }
            writeBytes(t, e = 0, r = 0) {
                0 === r && (r = t.length - e);
                let n = this.bufferU8.byteLength >>> 0;
                for (this._length < this.position + r && (this._length = this.position + r); n < this._length; ) n *= 2;
                if (n === this.bufferU8.byteLength) this.bufferU8.set(new Uint8ClampedArray(t.bufferU8.buffer, e, r), this.position);
                else {
                    const o = ArrayBufferPool.get(n);
                    o.set(this.bufferU8), o.set(new Uint8ClampedArray(t.bufferU8.buffer, e, r), this.position), ArrayBufferPool.add(this.bufferU8), (this.bufferU8 = o), (this.dataView = new DataView(this.bufferU8.buffer));
                }
                this.position += r;
            }
            readByte() {
                if (this.position + ByteArray.SIZE_OF_INT8 > this._length) throw new Error("Failed to read past end of the stream");
                const t = this.dataView.getInt8(this.position);
                return this._position++, t;
            }
            readUnsignedByte() {
                if (this.position + ByteArray.SIZE_OF_INT8 > this._length) throw new Error("Failed to read past end of the stream");
                const t = this.dataView.getUint8(this.position);
                return this._position++, t;
            }
            readIntLeb() {
                let result = (0, big_integer_1.default)(0),
                    value = (0, big_integer_1.default)(0),
                    shift = (0, big_integer_1.default)(0);
                do {
                    (value = (0, big_integer_1.default)(this.dataView.getUint8(this._position))), this._position++;
                    let t = value.and((0, big_integer_1.default)(127));
                    (t = t.shiftLeft((0, big_integer_1.default)(shift))), (result = result.or(t)), (shift = shift.plus((0, big_integer_1.default)(7)));
                } while (value.greaterOrEquals((0, big_integer_1.default)(128)));
                return result.greaterOrEquals(Number.MAX_SAFE_INTEGER) ? eval("String(result)") : result.toJSNumber();
            }
            readLebUTF() {
                return this.readUTFBytes(this.readIntLeb());
            }
            readFloat() {
                const t = this.dataView.getFloat32(this.position, !0);
                return (this._position += 4), t;
            }
            toString() {
                const t = this.fromCharCode();
                try {
                    return decodeURIComponent(escape(t));
                } catch (e) {
                    return t;
                }
            }
            ToString() {
                return "[BA: Length:" + this.length + "]";
            }
            clear() {
                (this.length = 0), (this.position = 0);
            }
            get(t) {
                return this.bufferU8[t];
            }
            set(t, e) {
                for (t === this._length && (this._length++, this._position++); this._length >= this.bufferU8.byteLength; ) this.expand();
                this.bufferU8[t] = e;
            }
            writeUTFLeb(t) {
                (t = unescape(encodeURIComponent(t))), this.writeIntLeb(t.length), this.writeUTFBytes(t);
            }
            writeByte(t) {
                for (this._length < this._position + ByteArray.SIZE_OF_INT8 && (this._length = this._position + ByteArray.SIZE_OF_INT8); this._length >= this.bufferU8.byteLength >>> 0; ) this.expand();
                this.dataView.setUint8(this.position, 255 & t), this.position++;
            }
            write(t) {
                for (this._length < this._position + (t >>> 0) && (this._length = this._position + (t >>> 0)); this._length >= this.bufferU8.byteLength; ) this.expand();
            }
            writeUTFBytes(t) {
                let e = 0;
                const r = t.length;
                for (; e < r - 10; e += 10)
                    this.writeByte(t.charCodeAt(e)),
                        this.writeByte(t.charCodeAt(e + 1)),
                        this.writeByte(t.charCodeAt(e + 2)),
                        this.writeByte(t.charCodeAt(e + 3)),
                        this.writeByte(t.charCodeAt(e + 4)),
                        this.writeByte(t.charCodeAt(e + 5)),
                        this.writeByte(t.charCodeAt(e + 6)),
                        this.writeByte(t.charCodeAt(e + 7)),
                        this.writeByte(t.charCodeAt(e + 8)),
                        this.writeByte(t.charCodeAt(e + 9));
                for (; e < r; e++) this.writeByte(t.charCodeAt(e));
            }
            readUTFBytes(t) {
                const e = ByteArray.STEP,
                    r = t + this.position;
                let n = null;
                for (let t = this.position; t < r; t += e) {
                    const o = Math.min(e, r - t),
                        i = String.fromCharCode.apply(String, new Uint8ClampedArray(this.bufferU8.buffer, t, o));
                    null == n ? (n = i) : (n += i);
                }
                null == n && (n = "");
                const o = decodeURIComponent(escape(n));
                return (this._position += t), o;
            }
            expand() {
                const t = ArrayBufferPool.get(2 * this.bufferU8.byteLength);
                t.set(this.bufferU8), ArrayBufferPool.add(this.bufferU8), (this.bufferU8 = t), (this.dataView = new DataView(this.bufferU8.buffer));
            }
            fromCharCode() {
                let t = "";
                const e = ByteArray.STEP;
                for (let r = 0; r < this.length; r += e) {
                    const n = Math.min(e, this.length - r);
                    t += String.fromCharCode.apply(String, new Uint8ClampedArray(this.bufferU8.buffer, r, n));
                }
                return t;
            }
            static encodeLeb128(t) {
                const e = [];
                let r = (0, big_integer_1.default)(t);
                do {
                    let t = r.and((0, big_integer_1.default)(127));
                    (r = r.shiftRight((0, big_integer_1.default)(7))), r.notEquals((0, big_integer_1.default)(0)) && (t = t.or((0, big_integer_1.default)(128))), e.push(Number(t));
                } while (r.notEquals(0));
                return e;
            }
        }
        (exports.ByteArray = ByteArray), (ByteArray.STEP = 65e3), (ByteArray.SIZE_OF_INT8 = 1);
    },
    function (t, e, r) {
        "use strict";
        function n(t) {
            return (n =
                "function" == typeof Symbol && "symbol" == typeof Symbol.iterator
                    ? function (t) {
                          return typeof t;
                      }
                    : function (t) {
                          return t && "function" == typeof Symbol && t.constructor === Symbol && t !== Symbol.prototype ? "symbol" : typeof t;
                      })(t);
        }
        function o(t, e, r) {
            var o = r.value;
            if ("function" != typeof o) throw new TypeError("@boundMethod decorator can only be applied to methods not: ".concat(n(o)));
            var i = !1;
            return {
                configurable: !0,
                get: function () {
                    if (i || this === t.prototype || this.hasOwnProperty(e) || "function" != typeof o) return o;
                    var r = o.bind(this);
                    return (
                        (i = !0),
                        Object.defineProperty(this, e, {
                            configurable: !0,
                            get: function () {
                                return r;
                            },
                            set: function (t) {
                                (o = t), delete this[e];
                            },
                        }),
                        (i = !1),
                        r
                    );
                },
                set: function (t) {
                    o = t;
                },
            };
        }
        function i(t) {
            var e;
            return (
                "undefined" != typeof Reflect && "function" == typeof Reflect.ownKeys
                    ? (e = Reflect.ownKeys(t.prototype))
                    : ((e = Object.getOwnPropertyNames(t.prototype)), "function" == typeof Object.getOwnPropertySymbols && (e = e.concat(Object.getOwnPropertySymbols(t.prototype)))),
                e.forEach(function (e) {
                    if ("constructor" !== e) {
                        var r = Object.getOwnPropertyDescriptor(t.prototype, e);
                        "function" == typeof r.value && Object.defineProperty(t.prototype, e, o(t, e, r));
                    }
                }),
                t
            );
        }
        function s() {
            return 1 === arguments.length ? i.apply(void 0, arguments) : o.apply(void 0, arguments);
        }
        r.r(e),
            r.d(e, "boundMethod", function () {
                return o;
            }),
            r.d(e, "boundClass", function () {
                return i;
            }),
            r.d(e, "default", function () {
                return s;
            });
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 });
        const n = r(3),
            o = r(4);
        self.onmessage = (t) => {
            switch (t.data.command) {
                case n.Command.OPEN_CONNECTION:
                    (r = t.data.message),
                        i
                            .connectTo(r.host, r.port, r.deviceType)
                            .then(() => {
                                self.postMessage({ command: n.Command.OPEN_CONNECTION, message: { status: !0 } });
                            })
                            .catch(() => {
                                self.postMessage({ command: n.Command.OPEN_CONNECTION, message: { status: !1 } });
                            });
                    break;
                case n.Command.SEND:
                    (e = t.data.message), i.sendData(e.type, ...e.args);
                    break;
                case n.Command.CLOSE_CONNECTION:
                    i.close();
            }
            var e, r;
        };
        const i = new o.Connection(
            (t) => {
                self.postMessage({ command: n.Command.READ, message: t });
            },
            function () {
                self.postMessage({ command: n.Command.CLOSE_CONNECTION, message: null });
            }.bind(self)
        );
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }),
            (e.Command = void 0),
            (function (t) {
                (t[(t.OPEN_CONNECTION = 0)] = "OPEN_CONNECTION"), (t[(t.CLOSE_CONNECTION = 1)] = "CLOSE_CONNECTION"), (t[(t.SEND = 2)] = "SEND"), (t[(t.READ = 3)] = "READ");
            })(e.Command || (e.Command = {}));
    },
    function (t, e, r) {
        "use strict";
        var n =
                (this && this.__decorate) ||
                function (t, e, r, n) {
                    var o,
                        i = arguments.length,
                        s = i < 3 ? e : null === n ? (n = Object.getOwnPropertyDescriptor(e, r)) : n;
                    if ("object" == typeof Reflect && "function" == typeof Reflect.decorate) s = Reflect.decorate(t, e, r, n);
                    else for (var a = t.length - 1; a >= 0; a--) (o = t[a]) && (s = (i < 3 ? o(s) : i > 3 ? o(e, r, s) : o(e, r)) || s);
                    return i > 3 && s && Object.defineProperty(e, r, s), s;
                },
            o =
                (this && this.__importDefault) ||
                function (t) {
                    return t && t.__esModule ? t : { default: t };
                };
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.Connection = void 0);
        const i = o(r(5)),
            s = r(1),
            a = r(0),
            u = r(14),
            l = r(16),
            I = r(17),
            c = r(19),
            p = r(20);
        class f extends u.Socket {
            constructor(t, e) {
                super(),
                    (this.messageHandler = t),
                    (this.closeHandler = e),
                    (this.clientPacketId = 0),
                    (this.serverPacketId = 0),
                    (this.deviceType = null),
                    (this.allowReconnect = !1),
                    this.onConnectionClosed.add(this.onCloseSocket),
                    this.onErrorOccured.add(this.onErrorSocket),
                    this.onDataReceived.add(this.onData);
            }
            sendData(t, ...e) {
                if (!this.connected) return void console.warn("Trying to send packet to closed connection");
                (0, i.default)(null != this.deviceType), console.log("Send packet " + p.ClientPacketType[t] + " type " + t + " id " + this.clientPacketId + " data: " + t, e);
                const r = new l.PacketClient(t, this.deviceType);
                r.load(e), this.writeIntLeb(r.length + a.ByteArray.encodeLeb128(this.clientPacketId).length), this.writeIntLeb(this.clientPacketId++), this.writeBytes(r), this.flush();
            }
            connectTo(t, e, r) {
                return console.log("Connecting to: " + t + ":" + e + " deviceType=" + r), (this.allowReconnect = !1), (this.deviceType = r), super.connect(t, e);
            }
            receiveData(t) {
                const e = t.readIntLeb(),
                    r = t.readIntLeb();
                if (r >= c.ServerPacketType.MAX_TYPE || r <= 0) throw new Error("Received server packet with wrong type " + r);
                const n = new I.PacketServer(r, e, t);
                if (e !== this.serverPacketId)
                    return (
                        console.warn(`Reading packet ${c.ServerPacketType[n.type]}:${n.type} with id ${n.id} but expected ${this.serverPacketId + 1} and length ${n.length}, data: ${n.bytesLength > 300 ? "..." : n}\n`),
                        console.warn("Error. Server packets was skipped"),
                        (this.allowReconnect = !0),
                        void this.close()
                    );
                (this.serverPacketId = e + 1), this.dispatchPacket(n);
            }
            dispatchPacket(t) {
                try {
                    this.messageHandler({ type: t.type, packet: t });
                } catch (t) {
                    console.log(t),
                        console.log("Stack:" + t.stack),
                        setTimeout(() => {
                            throw t;
                        }, 1);
                }
            }
            onErrorSocket() {
                console.log("onErrorSocket "), this.close(), (this.clientPacketId = 0), (this.serverPacketId = 0);
            }
            onCloseSocket() {
                (this.clientPacketId = 0), (this.serverPacketId = 0), this.closeHandler();
            }
            onData() {
                let t = new a.ByteArray();
                if (!this.connected || this.bytesAvailable < 4) return;
                let e = this.readIntLeb();
                e > this.bytesAvailable || (this.readBytes(t, 0, e), this.receiveData(t));
            }
        }
        n([s.boundMethod], f.prototype, "dispatchPacket", null),
            n([s.boundMethod], f.prototype, "onErrorSocket", null),
            n([s.boundMethod], f.prototype, "onCloseSocket", null),
            n([s.boundMethod], f.prototype, "onData", null),
            (e.Connection = f);
    },
    function (t, e, r) {
        "use strict";
        (function (e) {
            var n = r(7);
            /*!
             * The buffer module from node.js, for the browser.
             *
             * @author   Feross Aboukhadijeh <feross@feross.org> <http://feross.org>
             * @license  MIT
             */ function o(t, e) {
                if (t === e) return 0;
                for (var r = t.length, n = e.length, o = 0, i = Math.min(r, n); o < i; ++o)
                    if (t[o] !== e[o]) {
                        (r = t[o]), (n = e[o]);
                        break;
                    }
                return r < n ? -1 : n < r ? 1 : 0;
            }
            function i(t) {
                return e.Buffer && "function" == typeof e.Buffer.isBuffer ? e.Buffer.isBuffer(t) : !(null == t || !t._isBuffer);
            }
            var s = r(8),
                a = Object.prototype.hasOwnProperty,
                u = Array.prototype.slice,
                l = "foo" === function () {}.name;
            function I(t) {
                return Object.prototype.toString.call(t);
            }
            function c(t) {
                return !i(t) && "function" == typeof e.ArrayBuffer && ("function" == typeof ArrayBuffer.isView ? ArrayBuffer.isView(t) : !!t && (t instanceof DataView || !!(t.buffer && t.buffer instanceof ArrayBuffer)));
            }
            var p = (t.exports = S),
                f = /\s*function\s+([^\(\s]*)\s*/;
            function E(t) {
                if (s.isFunction(t)) {
                    if (l) return t.name;
                    var e = t.toString().match(f);
                    return e && e[1];
                }
            }
            function h(t, e) {
                return "string" == typeof t ? (t.length < e ? t : t.slice(0, e)) : t;
            }
            function y(t) {
                if (l || !s.isFunction(t)) return s.inspect(t);
                var e = E(t);
                return "[Function" + (e ? ": " + e : "") + "]";
            }
            function _(t, e, r, n, o) {
                throw new p.AssertionError({ message: r, actual: t, expected: e, operator: n, stackStartFunction: o });
            }
            function S(t, e) {
                t || _(t, !0, e, "==", p.ok);
            }
            function O(t, e, r, n) {
                if (t === e) return !0;
                if (i(t) && i(e)) return 0 === o(t, e);
                if (s.isDate(t) && s.isDate(e)) return t.getTime() === e.getTime();
                if (s.isRegExp(t) && s.isRegExp(e)) return t.source === e.source && t.global === e.global && t.multiline === e.multiline && t.lastIndex === e.lastIndex && t.ignoreCase === e.ignoreCase;
                if ((null !== t && "object" == typeof t) || (null !== e && "object" == typeof e)) {
                    if (c(t) && c(e) && I(t) === I(e) && !(t instanceof Float32Array || t instanceof Float64Array)) return 0 === o(new Uint8Array(t.buffer), new Uint8Array(e.buffer));
                    if (i(t) !== i(e)) return !1;
                    var a = (n = n || { actual: [], expected: [] }).actual.indexOf(t);
                    return (
                        (-1 !== a && a === n.expected.indexOf(e)) ||
                        (n.actual.push(t),
                        n.expected.push(e),
                        (function (t, e, r, n) {
                            if (null == t || null == e) return !1;
                            if (s.isPrimitive(t) || s.isPrimitive(e)) return t === e;
                            if (r && Object.getPrototypeOf(t) !== Object.getPrototypeOf(e)) return !1;
                            var o = T(t),
                                i = T(e);
                            if ((o && !i) || (!o && i)) return !1;
                            if (o) return (t = u.call(t)), (e = u.call(e)), O(t, e, r);
                            var a,
                                l,
                                I = A(t),
                                c = A(e);
                            if (I.length !== c.length) return !1;
                            for (I.sort(), c.sort(), l = I.length - 1; l >= 0; l--) if (I[l] !== c[l]) return !1;
                            for (l = I.length - 1; l >= 0; l--) if (((a = I[l]), !O(t[a], e[a], r, n))) return !1;
                            return !0;
                        })(t, e, r, n))
                    );
                }
                return r ? t === e : t == e;
            }
            function T(t) {
                return "[object Arguments]" == Object.prototype.toString.call(t);
            }
            function d(t, e) {
                if (!t || !e) return !1;
                if ("[object RegExp]" == Object.prototype.toString.call(e)) return e.test(t);
                try {
                    if (t instanceof e) return !0;
                } catch (t) {}
                return !Error.isPrototypeOf(e) && !0 === e.call({}, t);
            }
            function g(t, e, r, n) {
                var o;
                if ("function" != typeof e) throw new TypeError('"block" argument must be a function');
                "string" == typeof r && ((n = r), (r = null)),
                    (o = (function (t) {
                        var e;
                        try {
                            t();
                        } catch (t) {
                            e = t;
                        }
                        return e;
                    })(e)),
                    (n = (r && r.name ? " (" + r.name + ")." : ".") + (n ? " " + n : ".")),
                    t && !o && _(o, r, "Missing expected exception" + n);
                var i = "string" == typeof n,
                    a = !t && o && !r;
                if ((((!t && s.isError(o) && i && d(o, r)) || a) && _(o, r, "Got unwanted exception" + n), (t && o && r && !d(o, r)) || (!t && o))) throw o;
            }
            (p.AssertionError = function (t) {
                (this.name = "AssertionError"),
                    (this.actual = t.actual),
                    (this.expected = t.expected),
                    (this.operator = t.operator),
                    t.message
                        ? ((this.message = t.message), (this.generatedMessage = !1))
                        : ((this.message = (function (t) {
                              return h(y(t.actual), 128) + " " + t.operator + " " + h(y(t.expected), 128);
                          })(this)),
                          (this.generatedMessage = !0));
                var e = t.stackStartFunction || _;
                if (Error.captureStackTrace) Error.captureStackTrace(this, e);
                else {
                    var r = new Error();
                    if (r.stack) {
                        var n = r.stack,
                            o = E(e),
                            i = n.indexOf("\n" + o);
                        if (i >= 0) {
                            var s = n.indexOf("\n", i + 1);
                            n = n.substring(s + 1);
                        }
                        this.stack = n;
                    }
                }
            }),
                s.inherits(p.AssertionError, Error),
                (p.fail = _),
                (p.ok = S),
                (p.equal = function (t, e, r) {
                    t != e && _(t, e, r, "==", p.equal);
                }),
                (p.notEqual = function (t, e, r) {
                    t == e && _(t, e, r, "!=", p.notEqual);
                }),
                (p.deepEqual = function (t, e, r) {
                    O(t, e, !1) || _(t, e, r, "deepEqual", p.deepEqual);
                }),
                (p.deepStrictEqual = function (t, e, r) {
                    O(t, e, !0) || _(t, e, r, "deepStrictEqual", p.deepStrictEqual);
                }),
                (p.notDeepEqual = function (t, e, r) {
                    O(t, e, !1) && _(t, e, r, "notDeepEqual", p.notDeepEqual);
                }),
                (p.notDeepStrictEqual = function t(e, r, n) {
                    O(e, r, !0) && _(e, r, n, "notDeepStrictEqual", t);
                }),
                (p.strictEqual = function (t, e, r) {
                    t !== e && _(t, e, r, "===", p.strictEqual);
                }),
                (p.notStrictEqual = function (t, e, r) {
                    t === e && _(t, e, r, "!==", p.notStrictEqual);
                }),
                (p.throws = function (t, e, r) {
                    g(!0, t, e, r);
                }),
                (p.doesNotThrow = function (t, e, r) {
                    g(!1, t, e, r);
                }),
                (p.ifError = function (t) {
                    if (t) throw t;
                }),
                (p.strict = n(
                    function t(e, r) {
                        e || _(e, !0, r, "==", t);
                    },
                    p,
                    { equal: p.strictEqual, deepEqual: p.deepStrictEqual, notEqual: p.notStrictEqual, notDeepEqual: p.notDeepStrictEqual }
                )),
                (p.strict.strict = p.strict);
            var A =
                Object.keys ||
                function (t) {
                    var e = [];
                    for (var r in t) a.call(t, r) && e.push(r);
                    return e;
                };
        }.call(this, r(6)));
    },
    function (t, e) {
        var r;
        r = (function () {
            return this;
        })();
        try {
            r = r || new Function("return this")();
        } catch (t) {
            "object" == typeof window && (r = window);
        }
        t.exports = r;
    },
    function (t, e, r) {
        "use strict";
        /*
object-assign
(c) Sindre Sorhus
@license MIT
*/ var n = Object.getOwnPropertySymbols,
            o = Object.prototype.hasOwnProperty,
            i = Object.prototype.propertyIsEnumerable;
        function s(t) {
            if (null == t) throw new TypeError("Object.assign cannot be called with null or undefined");
            return Object(t);
        }
        t.exports = (function () {
            try {
                if (!Object.assign) return !1;
                var t = new String("abc");
                if (((t[5] = "de"), "5" === Object.getOwnPropertyNames(t)[0])) return !1;
                for (var e = {}, r = 0; r < 10; r++) e["_" + String.fromCharCode(r)] = r;
                if (
                    "0123456789" !==
                    Object.getOwnPropertyNames(e)
                        .map(function (t) {
                            return e[t];
                        })
                        .join("")
                )
                    return !1;
                var n = {};
                return (
                    "abcdefghijklmnopqrst".split("").forEach(function (t) {
                        n[t] = t;
                    }),
                    "abcdefghijklmnopqrst" === Object.keys(Object.assign({}, n)).join("")
                );
            } catch (t) {
                return !1;
            }
        })()
            ? Object.assign
            : function (t, e) {
                  for (var r, a, u = s(t), l = 1; l < arguments.length; l++) {
                      for (var I in (r = Object(arguments[l]))) o.call(r, I) && (u[I] = r[I]);
                      if (n) {
                          a = n(r);
                          for (var c = 0; c < a.length; c++) i.call(r, a[c]) && (u[a[c]] = r[a[c]]);
                      }
                  }
                  return u;
              };
    },
    function (t, e, r) {
        (function (t) {
            var n =
                    Object.getOwnPropertyDescriptors ||
                    function (t) {
                        for (var e = Object.keys(t), r = {}, n = 0; n < e.length; n++) r[e[n]] = Object.getOwnPropertyDescriptor(t, e[n]);
                        return r;
                    },
                o = /%[sdj%]/g;
            (e.format = function (t) {
                if (!_(t)) {
                    for (var e = [], r = 0; r < arguments.length; r++) e.push(a(arguments[r]));
                    return e.join(" ");
                }
                r = 1;
                for (
                    var n = arguments,
                        i = n.length,
                        s = String(t).replace(o, function (t) {
                            if ("%%" === t) return "%";
                            if (r >= i) return t;
                            switch (t) {
                                case "%s":
                                    return String(n[r++]);
                                case "%d":
                                    return Number(n[r++]);
                                case "%j":
                                    try {
                                        return JSON.stringify(n[r++]);
                                    } catch (t) {
                                        return "[Circular]";
                                    }
                                default:
                                    return t;
                            }
                        }),
                        u = n[r];
                    r < i;
                    u = n[++r]
                )
                    h(u) || !T(u) ? (s += " " + u) : (s += " " + a(u));
                return s;
            }),
                (e.deprecate = function (r, n) {
                    if (void 0 !== t && !0 === t.noDeprecation) return r;
                    if (void 0 === t)
                        return function () {
                            return e.deprecate(r, n).apply(this, arguments);
                        };
                    var o = !1;
                    return function () {
                        if (!o) {
                            if (t.throwDeprecation) throw new Error(n);
                            t.traceDeprecation ? console.trace(n) : console.error(n), (o = !0);
                        }
                        return r.apply(this, arguments);
                    };
                });
            var i,
                s = {};
            function a(t, r) {
                var n = { seen: [], stylize: l };
                return (
                    arguments.length >= 3 && (n.depth = arguments[2]),
                    arguments.length >= 4 && (n.colors = arguments[3]),
                    E(r) ? (n.showHidden = r) : r && e._extend(n, r),
                    S(n.showHidden) && (n.showHidden = !1),
                    S(n.depth) && (n.depth = 2),
                    S(n.colors) && (n.colors = !1),
                    S(n.customInspect) && (n.customInspect = !0),
                    n.colors && (n.stylize = u),
                    I(n, t, n.depth)
                );
            }
            function u(t, e) {
                var r = a.styles[e];
                return r ? "[" + a.colors[r][0] + "m" + t + "[" + a.colors[r][1] + "m" : t;
            }
            function l(t, e) {
                return t;
            }
            function I(t, r, n) {
                if (t.customInspect && r && A(r.inspect) && r.inspect !== e.inspect && (!r.constructor || r.constructor.prototype !== r)) {
                    var o = r.inspect(n, t);
                    return _(o) || (o = I(t, o, n)), o;
                }
                var i = (function (t, e) {
                    if (S(e)) return t.stylize("undefined", "undefined");
                    if (_(e)) {
                        var r = "'" + JSON.stringify(e).replace(/^"|"$/g, "").replace(/'/g, "\\'").replace(/\\"/g, '"') + "'";
                        return t.stylize(r, "string");
                    }
                    if (y(e)) return t.stylize("" + e, "number");
                    if (E(e)) return t.stylize("" + e, "boolean");
                    if (h(e)) return t.stylize("null", "null");
                })(t, r);
                if (i) return i;
                var s = Object.keys(r),
                    a = (function (t) {
                        var e = {};
                        return (
                            t.forEach(function (t, r) {
                                e[t] = !0;
                            }),
                            e
                        );
                    })(s);
                if ((t.showHidden && (s = Object.getOwnPropertyNames(r)), g(r) && (s.indexOf("message") >= 0 || s.indexOf("description") >= 0))) return c(r);
                if (0 === s.length) {
                    if (A(r)) {
                        var u = r.name ? ": " + r.name : "";
                        return t.stylize("[Function" + u + "]", "special");
                    }
                    if (O(r)) return t.stylize(RegExp.prototype.toString.call(r), "regexp");
                    if (d(r)) return t.stylize(Date.prototype.toString.call(r), "date");
                    if (g(r)) return c(r);
                }
                var l,
                    T = "",
                    B = !1,
                    v = ["{", "}"];
                (f(r) && ((B = !0), (v = ["[", "]"])), A(r)) && (T = " [Function" + (r.name ? ": " + r.name : "") + "]");
                return (
                    O(r) && (T = " " + RegExp.prototype.toString.call(r)),
                    d(r) && (T = " " + Date.prototype.toUTCString.call(r)),
                    g(r) && (T = " " + c(r)),
                    0 !== s.length || (B && 0 != r.length)
                        ? n < 0
                            ? O(r)
                                ? t.stylize(RegExp.prototype.toString.call(r), "regexp")
                                : t.stylize("[Object]", "special")
                            : (t.seen.push(r),
                              (l = B
                                  ? (function (t, e, r, n, o) {
                                        for (var i = [], s = 0, a = e.length; s < a; ++s) D(e, String(s)) ? i.push(p(t, e, r, n, String(s), !0)) : i.push("");
                                        return (
                                            o.forEach(function (o) {
                                                o.match(/^\d+$/) || i.push(p(t, e, r, n, o, !0));
                                            }),
                                            i
                                        );
                                    })(t, r, n, a, s)
                                  : s.map(function (e) {
                                        return p(t, r, n, a, e, B);
                                    })),
                              t.seen.pop(),
                              (function (t, e, r) {
                                  if (
                                      t.reduce(function (t, e) {
                                          return e.indexOf("\n") >= 0 && 0, t + e.replace(/\u001b\[\d\d?m/g, "").length + 1;
                                      }, 0) > 60
                                  )
                                      return r[0] + ("" === e ? "" : e + "\n ") + " " + t.join(",\n  ") + " " + r[1];
                                  return r[0] + e + " " + t.join(", ") + " " + r[1];
                              })(l, T, v))
                        : v[0] + T + v[1]
                );
            }
            function c(t) {
                return "[" + Error.prototype.toString.call(t) + "]";
            }
            function p(t, e, r, n, o, i) {
                var s, a, u;
                if (
                    ((u = Object.getOwnPropertyDescriptor(e, o) || { value: e[o] }).get ? (a = u.set ? t.stylize("[Getter/Setter]", "special") : t.stylize("[Getter]", "special")) : u.set && (a = t.stylize("[Setter]", "special")),
                    D(n, o) || (s = "[" + o + "]"),
                    a ||
                        (t.seen.indexOf(u.value) < 0
                            ? (a = h(r) ? I(t, u.value, null) : I(t, u.value, r - 1)).indexOf("\n") > -1 &&
                              (a = i
                                  ? a
                                        .split("\n")
                                        .map(function (t) {
                                            return "  " + t;
                                        })
                                        .join("\n")
                                        .substr(2)
                                  : "\n" +
                                    a
                                        .split("\n")
                                        .map(function (t) {
                                            return "   " + t;
                                        })
                                        .join("\n"))
                            : (a = t.stylize("[Circular]", "special"))),
                    S(s))
                ) {
                    if (i && o.match(/^\d+$/)) return a;
                    (s = JSON.stringify("" + o)).match(/^"([a-zA-Z_][a-zA-Z_0-9]*)"$/)
                        ? ((s = s.substr(1, s.length - 2)), (s = t.stylize(s, "name")))
                        : ((s = s
                              .replace(/'/g, "\\'")
                              .replace(/\\"/g, '"')
                              .replace(/(^"|"$)/g, "'")),
                          (s = t.stylize(s, "string")));
                }
                return s + ": " + a;
            }
            function f(t) {
                return Array.isArray(t);
            }
            function E(t) {
                return "boolean" == typeof t;
            }
            function h(t) {
                return null === t;
            }
            function y(t) {
                return "number" == typeof t;
            }
            function _(t) {
                return "string" == typeof t;
            }
            function S(t) {
                return void 0 === t;
            }
            function O(t) {
                return T(t) && "[object RegExp]" === B(t);
            }
            function T(t) {
                return "object" == typeof t && null !== t;
            }
            function d(t) {
                return T(t) && "[object Date]" === B(t);
            }
            function g(t) {
                return T(t) && ("[object Error]" === B(t) || t instanceof Error);
            }
            function A(t) {
                return "function" == typeof t;
            }
            function B(t) {
                return Object.prototype.toString.call(t);
            }
            function v(t) {
                return t < 10 ? "0" + t.toString(10) : t.toString(10);
            }
            (e.debuglog = function (r) {
                if ((S(i) && (i = t.env.NODE_DEBUG || ""), (r = r.toUpperCase()), !s[r]))
                    if (new RegExp("\\b" + r + "\\b", "i").test(i)) {
                        var n = t.pid;
                        s[r] = function () {
                            var t = e.format.apply(e, arguments);
                            console.error("%s %d: %s", r, n, t);
                        };
                    } else s[r] = function () {};
                return s[r];
            }),
                (e.inspect = a),
                (a.colors = {
                    bold: [1, 22],
                    italic: [3, 23],
                    underline: [4, 24],
                    inverse: [7, 27],
                    white: [37, 39],
                    grey: [90, 39],
                    black: [30, 39],
                    blue: [34, 39],
                    cyan: [36, 39],
                    green: [32, 39],
                    magenta: [35, 39],
                    red: [31, 39],
                    yellow: [33, 39],
                }),
                (a.styles = { special: "cyan", number: "yellow", boolean: "yellow", undefined: "grey", null: "bold", string: "green", date: "magenta", regexp: "red" }),
                (e.isArray = f),
                (e.isBoolean = E),
                (e.isNull = h),
                (e.isNullOrUndefined = function (t) {
                    return null == t;
                }),
                (e.isNumber = y),
                (e.isString = _),
                (e.isSymbol = function (t) {
                    return "symbol" == typeof t;
                }),
                (e.isUndefined = S),
                (e.isRegExp = O),
                (e.isObject = T),
                (e.isDate = d),
                (e.isError = g),
                (e.isFunction = A),
                (e.isPrimitive = function (t) {
                    return null === t || "boolean" == typeof t || "number" == typeof t || "string" == typeof t || "symbol" == typeof t || void 0 === t;
                }),
                (e.isBuffer = r(10));
            var N = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
            function b() {
                var t = new Date(),
                    e = [v(t.getHours()), v(t.getMinutes()), v(t.getSeconds())].join(":");
                return [t.getDate(), N[t.getMonth()], e].join(" ");
            }
            function D(t, e) {
                return Object.prototype.hasOwnProperty.call(t, e);
            }
            (e.log = function () {
                console.log("%s - %s", b(), e.format.apply(e, arguments));
            }),
                (e.inherits = r(11)),
                (e._extend = function (t, e) {
                    if (!e || !T(e)) return t;
                    for (var r = Object.keys(e), n = r.length; n--; ) t[r[n]] = e[r[n]];
                    return t;
                });
            var R = "undefined" != typeof Symbol ? Symbol("util.promisify.custom") : void 0;
            function m(t, e) {
                if (!t) {
                    var r = new Error("Promise was rejected with a falsy value");
                    (r.reason = t), (t = r);
                }
                return e(t);
            }
            (e.promisify = function (t) {
                if ("function" != typeof t) throw new TypeError('The "original" argument must be of type Function');
                if (R && t[R]) {
                    var e;
                    if ("function" != typeof (e = t[R])) throw new TypeError('The "util.promisify.custom" argument must be of type Function');
                    return Object.defineProperty(e, R, { value: e, enumerable: !1, writable: !1, configurable: !0 }), e;
                }
                function e() {
                    for (
                        var e,
                            r,
                            n = new Promise(function (t, n) {
                                (e = t), (r = n);
                            }),
                            o = [],
                            i = 0;
                        i < arguments.length;
                        i++
                    )
                        o.push(arguments[i]);
                    o.push(function (t, n) {
                        t ? r(t) : e(n);
                    });
                    try {
                        t.apply(this, o);
                    } catch (t) {
                        r(t);
                    }
                    return n;
                }
                return Object.setPrototypeOf(e, Object.getPrototypeOf(t)), R && Object.defineProperty(e, R, { value: e, enumerable: !1, writable: !1, configurable: !0 }), Object.defineProperties(e, n(t));
            }),
                (e.promisify.custom = R),
                (e.callbackify = function (e) {
                    if ("function" != typeof e) throw new TypeError('The "original" argument must be of type Function');
                    function r() {
                        for (var r = [], n = 0; n < arguments.length; n++) r.push(arguments[n]);
                        var o = r.pop();
                        if ("function" != typeof o) throw new TypeError("The last argument must be of type Function");
                        var i = this,
                            s = function () {
                                return o.apply(i, arguments);
                            };
                        e.apply(this, r).then(
                            function (e) {
                                t.nextTick(s, null, e);
                            },
                            function (e) {
                                t.nextTick(m, e, s);
                            }
                        );
                    }
                    return Object.setPrototypeOf(r, Object.getPrototypeOf(e)), Object.defineProperties(r, n(e)), r;
                });
        }.call(this, r(9)));
    },
    function (t, e) {
        var r,
            n,
            o = (t.exports = {});
        function i() {
            throw new Error("setTimeout has not been defined");
        }
        function s() {
            throw new Error("clearTimeout has not been defined");
        }
        function a(t) {
            if (r === setTimeout) return setTimeout(t, 0);
            if ((r === i || !r) && setTimeout) return (r = setTimeout), setTimeout(t, 0);
            try {
                return r(t, 0);
            } catch (e) {
                try {
                    return r.call(null, t, 0);
                } catch (e) {
                    return r.call(this, t, 0);
                }
            }
        }
        !(function () {
            try {
                r = "function" == typeof setTimeout ? setTimeout : i;
            } catch (t) {
                r = i;
            }
            try {
                n = "function" == typeof clearTimeout ? clearTimeout : s;
            } catch (t) {
                n = s;
            }
        })();
        var u,
            l = [],
            I = !1,
            c = -1;
        function p() {
            I && u && ((I = !1), u.length ? (l = u.concat(l)) : (c = -1), l.length && f());
        }
        function f() {
            if (!I) {
                var t = a(p);
                I = !0;
                for (var e = l.length; e; ) {
                    for (u = l, l = []; ++c < e; ) u && u[c].run();
                    (c = -1), (e = l.length);
                }
                (u = null),
                    (I = !1),
                    (function (t) {
                        if (n === clearTimeout) return clearTimeout(t);
                        if ((n === s || !n) && clearTimeout) return (n = clearTimeout), clearTimeout(t);
                        try {
                            n(t);
                        } catch (e) {
                            try {
                                return n.call(null, t);
                            } catch (e) {
                                return n.call(this, t);
                            }
                        }
                    })(t);
            }
        }
        function E(t, e) {
            (this.fun = t), (this.array = e);
        }
        function h() {}
        (o.nextTick = function (t) {
            var e = new Array(arguments.length - 1);
            if (arguments.length > 1) for (var r = 1; r < arguments.length; r++) e[r - 1] = arguments[r];
            l.push(new E(t, e)), 1 !== l.length || I || a(f);
        }),
            (E.prototype.run = function () {
                this.fun.apply(null, this.array);
            }),
            (o.title = "browser"),
            (o.browser = !0),
            (o.env = {}),
            (o.argv = []),
            (o.version = ""),
            (o.versions = {}),
            (o.on = h),
            (o.addListener = h),
            (o.once = h),
            (o.off = h),
            (o.removeListener = h),
            (o.removeAllListeners = h),
            (o.emit = h),
            (o.prependListener = h),
            (o.prependOnceListener = h),
            (o.listeners = function (t) {
                return [];
            }),
            (o.binding = function (t) {
                throw new Error("process.binding is not supported");
            }),
            (o.cwd = function () {
                return "/";
            }),
            (o.chdir = function (t) {
                throw new Error("process.chdir is not supported");
            }),
            (o.umask = function () {
                return 0;
            });
    },
    function (t, e) {
        t.exports = function (t) {
            return t && "object" == typeof t && "function" == typeof t.copy && "function" == typeof t.fill && "function" == typeof t.readUInt8;
        };
    },
    function (t, e) {
        "function" == typeof Object.create
            ? (t.exports = function (t, e) {
                  (t.super_ = e), (t.prototype = Object.create(e.prototype, { constructor: { value: t, enumerable: !1, writable: !0, configurable: !0 } }));
              })
            : (t.exports = function (t, e) {
                  t.super_ = e;
                  var r = function () {};
                  (r.prototype = e.prototype), (t.prototype = new r()), (t.prototype.constructor = t);
              });
    },
    function (t, e, r) {
        (function (t) {
            var n,
                o = (function (t) {
                    "use strict";
                    var e = 1e7,
                        r = 9007199254740992,
                        n = c(r),
                        i = "function" == typeof BigInt;
                    function s(t, e, r, n) {
                        return void 0 === t ? s[0] : void 0 !== e && (10 != +e || r) ? k(t, e, r, n) : q(t);
                    }
                    function a(t, e) {
                        (this.value = t), (this.sign = e), (this.isSmall = !1);
                    }
                    function u(t) {
                        (this.value = t), (this.sign = t < 0), (this.isSmall = !0);
                    }
                    function l(t) {
                        this.value = t;
                    }
                    function I(t) {
                        return -r < t && t < r;
                    }
                    function c(t) {
                        return t < 1e7 ? [t] : t < 1e14 ? [t % 1e7, Math.floor(t / 1e7)] : [t % 1e7, Math.floor(t / 1e7) % 1e7, Math.floor(t / 1e14)];
                    }
                    function p(t) {
                        f(t);
                        var r = t.length;
                        if (r < 4 && D(t, n) < 0)
                            switch (r) {
                                case 0:
                                    return 0;
                                case 1:
                                    return t[0];
                                case 2:
                                    return t[0] + t[1] * e;
                                default:
                                    return t[0] + (t[1] + t[2] * e) * e;
                            }
                        return t;
                    }
                    function f(t) {
                        for (var e = t.length; 0 === t[--e]; );
                        t.length = e + 1;
                    }
                    function E(t) {
                        for (var e = new Array(t), r = -1; ++r < t; ) e[r] = 0;
                        return e;
                    }
                    function h(t) {
                        return t > 0 ? Math.floor(t) : Math.ceil(t);
                    }
                    function y(t, r) {
                        var n,
                            o,
                            i = t.length,
                            s = r.length,
                            a = new Array(i),
                            u = 0,
                            l = e;
                        for (o = 0; o < s; o++) (u = (n = t[o] + r[o] + u) >= l ? 1 : 0), (a[o] = n - u * l);
                        for (; o < i; ) (u = (n = t[o] + u) === l ? 1 : 0), (a[o++] = n - u * l);
                        return u > 0 && a.push(u), a;
                    }
                    function _(t, e) {
                        return t.length >= e.length ? y(t, e) : y(e, t);
                    }
                    function S(t, r) {
                        var n,
                            o,
                            i = t.length,
                            s = new Array(i),
                            a = e;
                        for (o = 0; o < i; o++) (n = t[o] - a + r), (r = Math.floor(n / a)), (s[o] = n - r * a), (r += 1);
                        for (; r > 0; ) (s[o++] = r % a), (r = Math.floor(r / a));
                        return s;
                    }
                    function O(t, e) {
                        var r,
                            n,
                            o = t.length,
                            i = e.length,
                            s = new Array(o),
                            a = 0;
                        for (r = 0; r < i; r++) (n = t[r] - a - e[r]) < 0 ? ((n += 1e7), (a = 1)) : (a = 0), (s[r] = n);
                        for (r = i; r < o; r++) {
                            if (!((n = t[r] - a) < 0)) {
                                s[r++] = n;
                                break;
                            }
                            (n += 1e7), (s[r] = n);
                        }
                        for (; r < o; r++) s[r] = t[r];
                        return f(s), s;
                    }
                    function T(t, e, r) {
                        var n,
                            o,
                            i = t.length,
                            s = new Array(i),
                            l = -e;
                        for (n = 0; n < i; n++) (o = t[n] + l), (l = Math.floor(o / 1e7)), (o %= 1e7), (s[n] = o < 0 ? o + 1e7 : o);
                        return "number" == typeof (s = p(s)) ? (r && (s = -s), new u(s)) : new a(s, r);
                    }
                    function d(t, e) {
                        var r,
                            n,
                            o,
                            i,
                            s = t.length,
                            a = e.length,
                            u = E(s + a);
                        for (o = 0; o < s; ++o) {
                            i = t[o];
                            for (var l = 0; l < a; ++l) (r = i * e[l] + u[o + l]), (n = Math.floor(r / 1e7)), (u[o + l] = r - 1e7 * n), (u[o + l + 1] += n);
                        }
                        return f(u), u;
                    }
                    function g(t, r) {
                        var n,
                            o,
                            i = t.length,
                            s = new Array(i),
                            a = e,
                            u = 0;
                        for (o = 0; o < i; o++) (n = t[o] * r + u), (u = Math.floor(n / a)), (s[o] = n - u * a);
                        for (; u > 0; ) (s[o++] = u % a), (u = Math.floor(u / a));
                        return s;
                    }
                    function A(t, e) {
                        for (var r = []; e-- > 0; ) r.push(0);
                        return r.concat(t);
                    }
                    function B(t, r, n) {
                        return new a(t < e ? g(r, t) : d(r, c(t)), n);
                    }
                    function v(t) {
                        var e,
                            r,
                            n,
                            o,
                            i = t.length,
                            s = E(i + i);
                        for (n = 0; n < i; n++) {
                            r = 0 - (o = t[n]) * o;
                            for (var a = n; a < i; a++) (e = o * t[a] * 2 + s[n + a] + r), (r = Math.floor(e / 1e7)), (s[n + a] = e - 1e7 * r);
                            s[n + i] = r;
                        }
                        return f(s), s;
                    }
                    function N(t, e) {
                        var r,
                            n,
                            o,
                            i,
                            s = t.length,
                            a = E(s);
                        for (o = 0, r = s - 1; r >= 0; --r) (o = (i = 1e7 * o + t[r]) - (n = h(i / e)) * e), (a[r] = 0 | n);
                        return [a, 0 | o];
                    }
                    function b(t, r) {
                        var n,
                            o = q(r);
                        if (i) return [new l(t.value / o.value), new l(t.value % o.value)];
                        var I,
                            y = t.value,
                            _ = o.value;
                        if (0 === _) throw new Error("Cannot divide by zero");
                        if (t.isSmall) return o.isSmall ? [new u(h(y / _)), new u(y % _)] : [s[0], t];
                        if (o.isSmall) {
                            if (1 === _) return [t, s[0]];
                            if (-1 == _) return [t.negate(), s[0]];
                            var S = Math.abs(_);
                            if (S < e) {
                                I = p((n = N(y, S))[0]);
                                var T = n[1];
                                return t.sign && (T = -T), "number" == typeof I ? (t.sign !== o.sign && (I = -I), [new u(I), new u(T)]) : [new a(I, t.sign !== o.sign), new u(T)];
                            }
                            _ = c(S);
                        }
                        var d = D(y, _);
                        if (-1 === d) return [s[0], t];
                        if (0 === d) return [s[t.sign === o.sign ? 1 : -1], s[0]];
                        I = (n =
                            y.length + _.length <= 200
                                ? (function (t, r) {
                                      var n,
                                          o,
                                          i,
                                          s,
                                          a,
                                          u,
                                          l,
                                          I = t.length,
                                          c = r.length,
                                          f = e,
                                          h = E(r.length),
                                          y = r[c - 1],
                                          _ = Math.ceil(f / (2 * y)),
                                          S = g(t, _),
                                          O = g(r, _);
                                      for (S.length <= I && S.push(0), O.push(0), y = O[c - 1], o = I - c; o >= 0; o--) {
                                          for (n = f - 1, S[o + c] !== y && (n = Math.floor((S[o + c] * f + S[o + c - 1]) / y)), i = 0, s = 0, u = O.length, a = 0; a < u; a++)
                                              (i += n * O[a]), (l = Math.floor(i / f)), (s += S[o + a] - (i - l * f)), (i = l), s < 0 ? ((S[o + a] = s + f), (s = -1)) : ((S[o + a] = s), (s = 0));
                                          for (; 0 !== s; ) {
                                              for (n -= 1, i = 0, a = 0; a < u; a++) (i += S[o + a] - f + O[a]) < 0 ? ((S[o + a] = i + f), (i = 0)) : ((S[o + a] = i), (i = 1));
                                              s += i;
                                          }
                                          h[o] = n;
                                      }
                                      return (S = N(S, _)[0]), [p(h), p(S)];
                                  })(y, _)
                                : (function (t, e) {
                                      for (var r, n, o, i, s, a = t.length, u = e.length, l = [], I = []; a; )
                                          if ((I.unshift(t[--a]), f(I), D(I, e) < 0)) l.push(0);
                                          else {
                                              (o = 1e7 * I[(n = I.length) - 1] + I[n - 2]), (i = 1e7 * e[u - 1] + e[u - 2]), n > u && (o = 1e7 * (o + 1)), (r = Math.ceil(o / i));
                                              do {
                                                  if (D((s = g(e, r)), I) <= 0) break;
                                                  r--;
                                              } while (r);
                                              l.push(r), (I = O(I, s));
                                          }
                                      return l.reverse(), [p(l), p(I)];
                                  })(y, _))[0];
                        var A = t.sign !== o.sign,
                            B = n[1],
                            v = t.sign;
                        return "number" == typeof I ? (A && (I = -I), (I = new u(I))) : (I = new a(I, A)), "number" == typeof B ? (v && (B = -B), (B = new u(B))) : (B = new a(B, v)), [I, B];
                    }
                    function D(t, e) {
                        if (t.length !== e.length) return t.length > e.length ? 1 : -1;
                        for (var r = t.length - 1; r >= 0; r--) if (t[r] !== e[r]) return t[r] > e[r] ? 1 : -1;
                        return 0;
                    }
                    function R(t) {
                        var e = t.abs();
                        return !e.isUnit() && (!!(e.equals(2) || e.equals(3) || e.equals(5)) || (!(e.isEven() || e.isDivisibleBy(3) || e.isDivisibleBy(5)) && (!!e.lesser(49) || void 0)));
                    }
                    function m(t, e) {
                        for (var r, n, i, s = t.prev(), a = s, u = 0; a.isEven(); ) (a = a.divide(2)), u++;
                        t: for (n = 0; n < e.length; n++)
                            if (!t.lesser(e[n]) && !(i = o(e[n]).modPow(a, t)).isUnit() && !i.equals(s)) {
                                for (r = u - 1; 0 != r; r--) {
                                    if ((i = i.square().mod(t)).isUnit()) return !1;
                                    if (i.equals(s)) continue t;
                                }
                                return !1;
                            }
                        return !0;
                    }
                    (a.prototype = Object.create(s.prototype)),
                        (u.prototype = Object.create(s.prototype)),
                        (l.prototype = Object.create(s.prototype)),
                        (a.prototype.add = function (t) {
                            var e = q(t);
                            if (this.sign !== e.sign) return this.subtract(e.negate());
                            var r = this.value,
                                n = e.value;
                            return e.isSmall ? new a(S(r, Math.abs(n)), this.sign) : new a(_(r, n), this.sign);
                        }),
                        (a.prototype.plus = a.prototype.add),
                        (u.prototype.add = function (t) {
                            var e = q(t),
                                r = this.value;
                            if (r < 0 !== e.sign) return this.subtract(e.negate());
                            var n = e.value;
                            if (e.isSmall) {
                                if (I(r + n)) return new u(r + n);
                                n = c(Math.abs(n));
                            }
                            return new a(S(n, Math.abs(r)), r < 0);
                        }),
                        (u.prototype.plus = u.prototype.add),
                        (l.prototype.add = function (t) {
                            return new l(this.value + q(t).value);
                        }),
                        (l.prototype.plus = l.prototype.add),
                        (a.prototype.subtract = function (t) {
                            var e = q(t);
                            if (this.sign !== e.sign) return this.add(e.negate());
                            var r = this.value,
                                n = e.value;
                            return e.isSmall
                                ? T(r, Math.abs(n), this.sign)
                                : (function (t, e, r) {
                                      var n;
                                      return D(t, e) >= 0 ? (n = O(t, e)) : ((n = O(e, t)), (r = !r)), "number" == typeof (n = p(n)) ? (r && (n = -n), new u(n)) : new a(n, r);
                                  })(r, n, this.sign);
                        }),
                        (a.prototype.minus = a.prototype.subtract),
                        (u.prototype.subtract = function (t) {
                            var e = q(t),
                                r = this.value;
                            if (r < 0 !== e.sign) return this.add(e.negate());
                            var n = e.value;
                            return e.isSmall ? new u(r - n) : T(n, Math.abs(r), r >= 0);
                        }),
                        (u.prototype.minus = u.prototype.subtract),
                        (l.prototype.subtract = function (t) {
                            return new l(this.value - q(t).value);
                        }),
                        (l.prototype.minus = l.prototype.subtract),
                        (a.prototype.negate = function () {
                            return new a(this.value, !this.sign);
                        }),
                        (u.prototype.negate = function () {
                            var t = this.sign,
                                e = new u(-this.value);
                            return (e.sign = !t), e;
                        }),
                        (l.prototype.negate = function () {
                            return new l(-this.value);
                        }),
                        (a.prototype.abs = function () {
                            return new a(this.value, !1);
                        }),
                        (u.prototype.abs = function () {
                            return new u(Math.abs(this.value));
                        }),
                        (l.prototype.abs = function () {
                            return new l(this.value >= 0 ? this.value : -this.value);
                        }),
                        (a.prototype.multiply = function (t) {
                            var r,
                                n,
                                o,
                                i = q(t),
                                u = this.value,
                                l = i.value,
                                I = this.sign !== i.sign;
                            if (i.isSmall) {
                                if (0 === l) return s[0];
                                if (1 === l) return this;
                                if (-1 === l) return this.negate();
                                if ((r = Math.abs(l)) < e) return new a(g(u, r), I);
                                l = c(r);
                            }
                            return (
                                (n = u.length),
                                (o = l.length),
                                new a(
                                    -0.012 * n - 0.012 * o + 15e-6 * n * o > 0
                                        ? (function t(e, r) {
                                              var n = Math.max(e.length, r.length);
                                              if (n <= 30) return d(e, r);
                                              n = Math.ceil(n / 2);
                                              var o = e.slice(n),
                                                  i = e.slice(0, n),
                                                  s = r.slice(n),
                                                  a = r.slice(0, n),
                                                  u = t(i, a),
                                                  l = t(o, s),
                                                  I = t(_(i, o), _(a, s)),
                                                  c = _(_(u, A(O(O(I, u), l), n)), A(l, 2 * n));
                                              return f(c), c;
                                          })(u, l)
                                        : d(u, l),
                                    I
                                )
                            );
                        }),
                        (a.prototype.times = a.prototype.multiply),
                        (u.prototype._multiplyBySmall = function (t) {
                            return I(t.value * this.value) ? new u(t.value * this.value) : B(Math.abs(t.value), c(Math.abs(this.value)), this.sign !== t.sign);
                        }),
                        (a.prototype._multiplyBySmall = function (t) {
                            return 0 === t.value ? s[0] : 1 === t.value ? this : -1 === t.value ? this.negate() : B(Math.abs(t.value), this.value, this.sign !== t.sign);
                        }),
                        (u.prototype.multiply = function (t) {
                            return q(t)._multiplyBySmall(this);
                        }),
                        (u.prototype.times = u.prototype.multiply),
                        (l.prototype.multiply = function (t) {
                            return new l(this.value * q(t).value);
                        }),
                        (l.prototype.times = l.prototype.multiply),
                        (a.prototype.square = function () {
                            return new a(v(this.value), !1);
                        }),
                        (u.prototype.square = function () {
                            var t = this.value * this.value;
                            return I(t) ? new u(t) : new a(v(c(Math.abs(this.value))), !1);
                        }),
                        (l.prototype.square = function (t) {
                            return new l(this.value * this.value);
                        }),
                        (a.prototype.divmod = function (t) {
                            var e = b(this, t);
                            return { quotient: e[0], remainder: e[1] };
                        }),
                        (l.prototype.divmod = u.prototype.divmod = a.prototype.divmod),
                        (a.prototype.divide = function (t) {
                            return b(this, t)[0];
                        }),
                        (l.prototype.over = l.prototype.divide = function (t) {
                            return new l(this.value / q(t).value);
                        }),
                        (u.prototype.over = u.prototype.divide = a.prototype.over = a.prototype.divide),
                        (a.prototype.mod = function (t) {
                            return b(this, t)[1];
                        }),
                        (l.prototype.mod = l.prototype.remainder = function (t) {
                            return new l(this.value % q(t).value);
                        }),
                        (u.prototype.remainder = u.prototype.mod = a.prototype.remainder = a.prototype.mod),
                        (a.prototype.pow = function (t) {
                            var e,
                                r,
                                n,
                                o = q(t),
                                i = this.value,
                                a = o.value;
                            if (0 === a) return s[1];
                            if (0 === i) return s[0];
                            if (1 === i) return s[1];
                            if (-1 === i) return o.isEven() ? s[1] : s[-1];
                            if (o.sign) return s[0];
                            if (!o.isSmall) throw new Error("The exponent " + o.toString() + " is too large.");
                            if (this.isSmall && I((e = Math.pow(i, a)))) return new u(h(e));
                            for (r = this, n = s[1]; !0 & a && ((n = n.times(r)), --a), 0 !== a; ) (a /= 2), (r = r.square());
                            return n;
                        }),
                        (u.prototype.pow = a.prototype.pow),
                        (l.prototype.pow = function (t) {
                            var e = q(t),
                                r = this.value,
                                n = e.value,
                                o = BigInt(0),
                                i = BigInt(1),
                                a = BigInt(2);
                            if (n === o) return s[1];
                            if (r === o) return s[0];
                            if (r === i) return s[1];
                            if (r === BigInt(-1)) return e.isEven() ? s[1] : s[-1];
                            if (e.isNegative()) return new l(o);
                            for (var u = this, I = s[1]; (n & i) === i && ((I = I.times(u)), --n), n !== o; ) (n /= a), (u = u.square());
                            return I;
                        }),
                        (a.prototype.modPow = function (t, e) {
                            if (((t = q(t)), (e = q(e)).isZero())) throw new Error("Cannot take modPow with modulus 0");
                            var r = s[1],
                                n = this.mod(e);
                            for (t.isNegative() && ((t = t.multiply(s[-1])), (n = n.modInv(e))); t.isPositive(); ) {
                                if (n.isZero()) return s[0];
                                t.isOdd() && (r = r.multiply(n).mod(e)), (t = t.divide(2)), (n = n.square().mod(e));
                            }
                            return r;
                        }),
                        (l.prototype.modPow = u.prototype.modPow = a.prototype.modPow),
                        (a.prototype.compareAbs = function (t) {
                            var e = q(t),
                                r = this.value,
                                n = e.value;
                            return e.isSmall ? 1 : D(r, n);
                        }),
                        (u.prototype.compareAbs = function (t) {
                            var e = q(t),
                                r = Math.abs(this.value),
                                n = e.value;
                            return e.isSmall ? (r === (n = Math.abs(n)) ? 0 : r > n ? 1 : -1) : -1;
                        }),
                        (l.prototype.compareAbs = function (t) {
                            var e = this.value,
                                r = q(t).value;
                            return (e = e >= 0 ? e : -e) === (r = r >= 0 ? r : -r) ? 0 : e > r ? 1 : -1;
                        }),
                        (a.prototype.compare = function (t) {
                            if (t === 1 / 0) return -1;
                            if (t === -1 / 0) return 1;
                            var e = q(t),
                                r = this.value,
                                n = e.value;
                            return this.sign !== e.sign ? (e.sign ? 1 : -1) : e.isSmall ? (this.sign ? -1 : 1) : D(r, n) * (this.sign ? -1 : 1);
                        }),
                        (a.prototype.compareTo = a.prototype.compare),
                        (u.prototype.compare = function (t) {
                            if (t === 1 / 0) return -1;
                            if (t === -1 / 0) return 1;
                            var e = q(t),
                                r = this.value,
                                n = e.value;
                            return e.isSmall ? (r == n ? 0 : r > n ? 1 : -1) : r < 0 !== e.sign ? (r < 0 ? -1 : 1) : r < 0 ? 1 : -1;
                        }),
                        (u.prototype.compareTo = u.prototype.compare),
                        (l.prototype.compare = function (t) {
                            if (t === 1 / 0) return -1;
                            if (t === -1 / 0) return 1;
                            var e = this.value,
                                r = q(t).value;
                            return e === r ? 0 : e > r ? 1 : -1;
                        }),
                        (l.prototype.compareTo = l.prototype.compare),
                        (a.prototype.equals = function (t) {
                            return 0 === this.compare(t);
                        }),
                        (l.prototype.eq = l.prototype.equals = u.prototype.eq = u.prototype.equals = a.prototype.eq = a.prototype.equals),
                        (a.prototype.notEquals = function (t) {
                            return 0 !== this.compare(t);
                        }),
                        (l.prototype.neq = l.prototype.notEquals = u.prototype.neq = u.prototype.notEquals = a.prototype.neq = a.prototype.notEquals),
                        (a.prototype.greater = function (t) {
                            return this.compare(t) > 0;
                        }),
                        (l.prototype.gt = l.prototype.greater = u.prototype.gt = u.prototype.greater = a.prototype.gt = a.prototype.greater),
                        (a.prototype.lesser = function (t) {
                            return this.compare(t) < 0;
                        }),
                        (l.prototype.lt = l.prototype.lesser = u.prototype.lt = u.prototype.lesser = a.prototype.lt = a.prototype.lesser),
                        (a.prototype.greaterOrEquals = function (t) {
                            return this.compare(t) >= 0;
                        }),
                        (l.prototype.geq = l.prototype.greaterOrEquals = u.prototype.geq = u.prototype.greaterOrEquals = a.prototype.geq = a.prototype.greaterOrEquals),
                        (a.prototype.lesserOrEquals = function (t) {
                            return this.compare(t) <= 0;
                        }),
                        (l.prototype.leq = l.prototype.lesserOrEquals = u.prototype.leq = u.prototype.lesserOrEquals = a.prototype.leq = a.prototype.lesserOrEquals),
                        (a.prototype.isEven = function () {
                            return 0 == (1 & this.value[0]);
                        }),
                        (u.prototype.isEven = function () {
                            return 0 == (1 & this.value);
                        }),
                        (l.prototype.isEven = function () {
                            return (this.value & BigInt(1)) === BigInt(0);
                        }),
                        (a.prototype.isOdd = function () {
                            return 1 == (1 & this.value[0]);
                        }),
                        (u.prototype.isOdd = function () {
                            return 1 == (1 & this.value);
                        }),
                        (l.prototype.isOdd = function () {
                            return (this.value & BigInt(1)) === BigInt(1);
                        }),
                        (a.prototype.isPositive = function () {
                            return !this.sign;
                        }),
                        (u.prototype.isPositive = function () {
                            return this.value > 0;
                        }),
                        (l.prototype.isPositive = u.prototype.isPositive),
                        (a.prototype.isNegative = function () {
                            return this.sign;
                        }),
                        (u.prototype.isNegative = function () {
                            return this.value < 0;
                        }),
                        (l.prototype.isNegative = u.prototype.isNegative),
                        (a.prototype.isUnit = function () {
                            return !1;
                        }),
                        (u.prototype.isUnit = function () {
                            return 1 === Math.abs(this.value);
                        }),
                        (l.prototype.isUnit = function () {
                            return this.abs().value === BigInt(1);
                        }),
                        (a.prototype.isZero = function () {
                            return !1;
                        }),
                        (u.prototype.isZero = function () {
                            return 0 === this.value;
                        }),
                        (l.prototype.isZero = function () {
                            return this.value === BigInt(0);
                        }),
                        (a.prototype.isDivisibleBy = function (t) {
                            var e = q(t);
                            return !e.isZero() && (!!e.isUnit() || (0 === e.compareAbs(2) ? this.isEven() : this.mod(e).isZero()));
                        }),
                        (l.prototype.isDivisibleBy = u.prototype.isDivisibleBy = a.prototype.isDivisibleBy),
                        (a.prototype.isPrime = function (t) {
                            var e = R(this);
                            if (void 0 !== e) return e;
                            var r = this.abs(),
                                n = r.bitLength();
                            if (n <= 64) return m(r, [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37]);
                            for (var i = Math.log(2) * n.toJSNumber(), s = Math.ceil(!0 === t ? 2 * Math.pow(i, 2) : i), a = [], u = 0; u < s; u++) a.push(o(u + 2));
                            return m(r, a);
                        }),
                        (l.prototype.isPrime = u.prototype.isPrime = a.prototype.isPrime),
                        (a.prototype.isProbablePrime = function (t, e) {
                            var r = R(this);
                            if (void 0 !== r) return r;
                            for (var n = this.abs(), i = void 0 === t ? 5 : t, s = [], a = 0; a < i; a++) s.push(o.randBetween(2, n.minus(2), e));
                            return m(n, s);
                        }),
                        (l.prototype.isProbablePrime = u.prototype.isProbablePrime = a.prototype.isProbablePrime),
                        (a.prototype.modInv = function (t) {
                            for (var e, r, n, i = o.zero, s = o.one, a = q(t), u = this.abs(); !u.isZero(); ) (e = a.divide(u)), (r = i), (n = a), (i = s), (a = u), (s = r.subtract(e.multiply(s))), (u = n.subtract(e.multiply(u)));
                            if (!a.isUnit()) throw new Error(this.toString() + " and " + t.toString() + " are not co-prime");
                            return -1 === i.compare(0) && (i = i.add(t)), this.isNegative() ? i.negate() : i;
                        }),
                        (l.prototype.modInv = u.prototype.modInv = a.prototype.modInv),
                        (a.prototype.next = function () {
                            var t = this.value;
                            return this.sign ? T(t, 1, this.sign) : new a(S(t, 1), this.sign);
                        }),
                        (u.prototype.next = function () {
                            var t = this.value;
                            return t + 1 < r ? new u(t + 1) : new a(n, !1);
                        }),
                        (l.prototype.next = function () {
                            return new l(this.value + BigInt(1));
                        }),
                        (a.prototype.prev = function () {
                            var t = this.value;
                            return this.sign ? new a(S(t, 1), !0) : T(t, 1, this.sign);
                        }),
                        (u.prototype.prev = function () {
                            var t = this.value;
                            return t - 1 > -r ? new u(t - 1) : new a(n, !0);
                        }),
                        (l.prototype.prev = function () {
                            return new l(this.value - BigInt(1));
                        });
                    for (var w = [1]; 2 * w[w.length - 1] <= e; ) w.push(2 * w[w.length - 1]);
                    var L = w.length,
                        P = w[L - 1];
                    function C(t) {
                        return Math.abs(t) <= e;
                    }
                    function G(t, e, r) {
                        e = q(e);
                        for (var n = t.isNegative(), i = e.isNegative(), s = n ? t.not() : t, a = i ? e.not() : e, u = 0, l = 0, I = null, c = null, p = []; !s.isZero() || !a.isZero(); )
                            (u = (I = b(s, P))[1].toJSNumber()), n && (u = P - 1 - u), (l = (c = b(a, P))[1].toJSNumber()), i && (l = P - 1 - l), (s = I[0]), (a = c[0]), p.push(r(u, l));
                        for (var f = 0 !== r(n ? 1 : 0, i ? 1 : 0) ? o(-1) : o(0), E = p.length - 1; E >= 0; E -= 1) f = f.multiply(P).add(o(p[E]));
                        return f;
                    }
                    (a.prototype.shiftLeft = function (t) {
                        var e = q(t).toJSNumber();
                        if (!C(e)) throw new Error(String(e) + " is too large for shifting.");
                        if (e < 0) return this.shiftRight(-e);
                        var r = this;
                        if (r.isZero()) return r;
                        for (; e >= L; ) (r = r.multiply(P)), (e -= L - 1);
                        return r.multiply(w[e]);
                    }),
                        (l.prototype.shiftLeft = u.prototype.shiftLeft = a.prototype.shiftLeft),
                        (a.prototype.shiftRight = function (t) {
                            var e,
                                r = q(t).toJSNumber();
                            if (!C(r)) throw new Error(String(r) + " is too large for shifting.");
                            if (r < 0) return this.shiftLeft(-r);
                            for (var n = this; r >= L; ) {
                                if (n.isZero() || (n.isNegative() && n.isUnit())) return n;
                                (n = (e = b(n, P))[1].isNegative() ? e[0].prev() : e[0]), (r -= L - 1);
                            }
                            return (e = b(n, w[r]))[1].isNegative() ? e[0].prev() : e[0];
                        }),
                        (l.prototype.shiftRight = u.prototype.shiftRight = a.prototype.shiftRight),
                        (a.prototype.not = function () {
                            return this.negate().prev();
                        }),
                        (l.prototype.not = u.prototype.not = a.prototype.not),
                        (a.prototype.and = function (t) {
                            return G(this, t, function (t, e) {
                                return t & e;
                            });
                        }),
                        (l.prototype.and = u.prototype.and = a.prototype.and),
                        (a.prototype.or = function (t) {
                            return G(this, t, function (t, e) {
                                return t | e;
                            });
                        }),
                        (l.prototype.or = u.prototype.or = a.prototype.or),
                        (a.prototype.xor = function (t) {
                            return G(this, t, function (t, e) {
                                return t ^ e;
                            });
                        }),
                        (l.prototype.xor = u.prototype.xor = a.prototype.xor);
                    function M(t) {
                        var r = t.value,
                            n = "number" == typeof r ? r | (1 << 30) : "bigint" == typeof r ? r | BigInt(1 << 30) : (r[0] + r[1] * e) | 1073758208;
                        return n & -n;
                    }
                    function U(t, e) {
                        return (t = q(t)), (e = q(e)), t.greater(e) ? t : e;
                    }
                    function W(t, e) {
                        return (t = q(t)), (e = q(e)), t.lesser(e) ? t : e;
                    }
                    function F(t, e) {
                        if (((t = q(t).abs()), (e = q(e).abs()), t.equals(e))) return t;
                        if (t.isZero()) return e;
                        if (e.isZero()) return t;
                        for (var r, n, o = s[1]; t.isEven() && e.isEven(); ) (r = W(M(t), M(e))), (t = t.divide(r)), (e = e.divide(r)), (o = o.multiply(r));
                        for (; t.isEven(); ) t = t.divide(M(t));
                        do {
                            for (; e.isEven(); ) e = e.divide(M(e));
                            t.greater(e) && ((n = e), (e = t), (t = n)), (e = e.subtract(t));
                        } while (!e.isZero());
                        return o.isUnit() ? t : t.multiply(o);
                    }
                    (a.prototype.bitLength = function () {
                        var t = this;
                        return (
                            t.compareTo(o(0)) < 0 && (t = t.negate().subtract(o(1))),
                            0 === t.compareTo(o(0))
                                ? o(0)
                                : o(
                                      (function t(e, r) {
                                          if (r.compareTo(e) <= 0) {
                                              var n = t(e, r.square(r)),
                                                  i = n.p,
                                                  s = n.e,
                                                  a = i.multiply(r);
                                              return a.compareTo(e) <= 0 ? { p: a, e: 2 * s + 1 } : { p: i, e: 2 * s };
                                          }
                                          return { p: o(1), e: 0 };
                                      })(t, o(2)).e
                                  ).add(o(1))
                        );
                    }),
                        (l.prototype.bitLength = u.prototype.bitLength = a.prototype.bitLength);
                    var k = function (t, e, r, n) {
                        (r = r || "0123456789abcdefghijklmnopqrstuvwxyz"), (t = String(t)), n || ((t = t.toLowerCase()), (r = r.toLowerCase()));
                        var o,
                            i = t.length,
                            s = Math.abs(e),
                            a = {};
                        for (o = 0; o < r.length; o++) a[r[o]] = o;
                        for (o = 0; o < i; o++) {
                            if ("-" !== (I = t[o]) && I in a && a[I] >= s) {
                                if ("1" === I && 1 === s) continue;
                                throw new Error(I + " is not a valid digit in base " + e + ".");
                            }
                        }
                        e = q(e);
                        var u = [],
                            l = "-" === t[0];
                        for (o = l ? 1 : 0; o < t.length; o++) {
                            var I;
                            if ((I = t[o]) in a) u.push(q(a[I]));
                            else {
                                if ("<" !== I) throw new Error(I + " is not a valid character");
                                var c = o;
                                do {
                                    o++;
                                } while (">" !== t[o] && o < t.length);
                                u.push(q(t.slice(c + 1, o)));
                            }
                        }
                        return V(u, e, l);
                    };
                    function V(t, e, r) {
                        var n,
                            o = s[0],
                            i = s[1];
                        for (n = t.length - 1; n >= 0; n--) (o = o.add(t[n].times(i))), (i = i.times(e));
                        return r ? o.negate() : o;
                    }
                    function H(t, e) {
                        if ((e = o(e)).isZero()) {
                            if (t.isZero()) return { value: [0], isNegative: !1 };
                            throw new Error("Cannot convert nonzero numbers to base 0.");
                        }
                        if (e.equals(-1)) {
                            if (t.isZero()) return { value: [0], isNegative: !1 };
                            if (t.isNegative()) return { value: [].concat.apply([], Array.apply(null, Array(-t.toJSNumber())).map(Array.prototype.valueOf, [1, 0])), isNegative: !1 };
                            var r = Array.apply(null, Array(t.toJSNumber() - 1)).map(Array.prototype.valueOf, [0, 1]);
                            return r.unshift([1]), { value: [].concat.apply([], r), isNegative: !1 };
                        }
                        var n = !1;
                        if ((t.isNegative() && e.isPositive() && ((n = !0), (t = t.abs())), e.isUnit()))
                            return t.isZero() ? { value: [0], isNegative: !1 } : { value: Array.apply(null, Array(t.toJSNumber())).map(Number.prototype.valueOf, 1), isNegative: n };
                        for (var i, s = [], a = t; a.isNegative() || a.compareAbs(e) >= 0; ) {
                            (i = a.divmod(e)), (a = i.quotient);
                            var u = i.remainder;
                            u.isNegative() && ((u = e.minus(u).abs()), (a = a.next())), s.push(u.toJSNumber());
                        }
                        return s.push(a.toJSNumber()), { value: s.reverse(), isNegative: n };
                    }
                    function j(t, e, r) {
                        var n = H(t, e);
                        return (
                            (n.isNegative ? "-" : "") +
                            n.value
                                .map(function (t) {
                                    return (function (t, e) {
                                        return t < (e = e || "0123456789abcdefghijklmnopqrstuvwxyz").length ? e[t] : "<" + t + ">";
                                    })(t, r);
                                })
                                .join("")
                        );
                    }
                    function x(t) {
                        if (I(+t)) {
                            var e = +t;
                            if (e === h(e)) return i ? new l(BigInt(e)) : new u(e);
                            throw new Error("Invalid integer: " + t);
                        }
                        var r = "-" === t[0];
                        r && (t = t.slice(1));
                        var n = t.split(/e/i);
                        if (n.length > 2) throw new Error("Invalid integer: " + n.join("e"));
                        if (2 === n.length) {
                            var o = n[1];
                            if (("+" === o[0] && (o = o.slice(1)), (o = +o) !== h(o) || !I(o))) throw new Error("Invalid integer: " + o + " is not a valid exponent.");
                            var s = n[0],
                                c = s.indexOf(".");
                            if ((c >= 0 && ((o -= s.length - c - 1), (s = s.slice(0, c) + s.slice(c + 1))), o < 0)) throw new Error("Cannot include negative exponent part for integers");
                            t = s += new Array(o + 1).join("0");
                        }
                        if (!/^([0-9][0-9]*)$/.test(t)) throw new Error("Invalid integer: " + t);
                        if (i) return new l(BigInt(r ? "-" + t : t));
                        for (var p = [], E = t.length, y = E - 7; E > 0; ) p.push(+t.slice(y, E)), (y -= 7) < 0 && (y = 0), (E -= 7);
                        return f(p), new a(p, r);
                    }
                    function q(t) {
                        return "number" == typeof t
                            ? (function (t) {
                                  if (i) return new l(BigInt(t));
                                  if (I(t)) {
                                      if (t !== h(t)) throw new Error(t + " is not an integer.");
                                      return new u(t);
                                  }
                                  return x(t.toString());
                              })(t)
                            : "string" == typeof t
                            ? x(t)
                            : "bigint" == typeof t
                            ? new l(t)
                            : t;
                    }
                    (a.prototype.toArray = function (t) {
                        return H(this, t);
                    }),
                        (u.prototype.toArray = function (t) {
                            return H(this, t);
                        }),
                        (l.prototype.toArray = function (t) {
                            return H(this, t);
                        }),
                        (a.prototype.toString = function (t, e) {
                            if ((void 0 === t && (t = 10), 10 !== t)) return j(this, t, e);
                            for (var r, n = this.value, o = n.length, i = String(n[--o]); --o >= 0; ) (r = String(n[o])), (i += "0000000".slice(r.length) + r);
                            return (this.sign ? "-" : "") + i;
                        }),
                        (u.prototype.toString = function (t, e) {
                            return void 0 === t && (t = 10), 10 != t ? j(this, t, e) : String(this.value);
                        }),
                        (l.prototype.toString = u.prototype.toString),
                        (l.prototype.toJSON = a.prototype.toJSON = u.prototype.toJSON = function () {
                            return this.toString();
                        }),
                        (a.prototype.valueOf = function () {
                            return parseInt(this.toString(), 10);
                        }),
                        (a.prototype.toJSNumber = a.prototype.valueOf),
                        (u.prototype.valueOf = function () {
                            return this.value;
                        }),
                        (u.prototype.toJSNumber = u.prototype.valueOf),
                        (l.prototype.valueOf = l.prototype.toJSNumber = function () {
                            return parseInt(this.toString(), 10);
                        });
                    for (var Y = 0; Y < 1e3; Y++) (s[Y] = q(Y)), Y > 0 && (s[-Y] = q(-Y));
                    return (
                        (s.one = s[1]),
                        (s.zero = s[0]),
                        (s.minusOne = s[-1]),
                        (s.max = U),
                        (s.min = W),
                        (s.gcd = F),
                        (s.lcm = function (t, e) {
                            return (t = q(t).abs()), (e = q(e).abs()), t.divide(F(t, e)).multiply(e);
                        }),
                        (s.isInstance = function (t) {
                            return t instanceof a || t instanceof u || t instanceof l;
                        }),
                        (s.randBetween = function (t, r, n) {
                            (t = q(t)), (r = q(r));
                            var o = n || Math.random,
                                i = W(t, r),
                                a = U(t, r).subtract(i).add(1);
                            if (a.isSmall) return i.add(Math.floor(o() * a));
                            for (var u = H(a, e).value, l = [], I = !0, c = 0; c < u.length; c++) {
                                var p = I ? u[c] + (c + 1 < u.length ? u[c + 1] / e : 0) : e,
                                    f = h(o() * p);
                                l.push(f), f < u[c] && (I = !1);
                            }
                            return i.add(s.fromArray(l, e, !1));
                        }),
                        (s.fromArray = function (t, e, r) {
                            return V(t.map(q), q(e || 10), r);
                        }),
                        s
                    );
                })();
            t.hasOwnProperty("exports") && (t.exports = o),
                void 0 ===
                    (n = function () {
                        return o;
                    }.call(e, r, e, t)) || (t.exports = n);
        }.call(this, r(13)(t)));
    },
    function (t, e) {
        t.exports = function (t) {
            return (
                t.webpackPolyfill ||
                    ((t.deprecate = function () {}),
                    (t.paths = []),
                    t.children || (t.children = []),
                    Object.defineProperty(t, "loaded", {
                        enumerable: !0,
                        get: function () {
                            return t.l;
                        },
                    }),
                    Object.defineProperty(t, "id", {
                        enumerable: !0,
                        get: function () {
                            return t.i;
                        },
                    }),
                    (t.webpackPolyfill = 1)),
                t
            );
        };
    },
    function (t, e, r) {
        "use strict";
        var n =
            (this && this.__decorate) ||
            function (t, e, r, n) {
                var o,
                    i = arguments.length,
                    s = i < 3 ? e : null === n ? (n = Object.getOwnPropertyDescriptor(e, r)) : n;
                if ("object" == typeof Reflect && "function" == typeof Reflect.decorate) s = Reflect.decorate(t, e, r, n);
                else for (var a = t.length - 1; a >= 0; a--) (o = t[a]) && (s = (i < 3 ? o(s) : i > 3 ? o(e, r, s) : o(e, r)) || s);
                return i > 3 && s && Object.defineProperty(e, r, s), s;
            };
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.Socket = void 0);
        const o = r(1),
            i = r(15),
            s = r(0);
        class a {
            constructor() {
                (this.client = null),
                    (this.inputByteArray = new s.ByteArray()),
                    (this.outputByteArray = new s.ByteArray()),
                    (this.onConnectionClosed = new i.SlotVoid()),
                    (this.onErrorOccured = new i.SlotVoid()),
                    (this.onDataReceived = new i.SlotVoid());
            }
            get isActive() {
                return null != this.client && this.client.readyState === this.client.OPEN;
            }
            get bytesAvailable() {
                return this.inputByteArray.length - this.inputByteArray.position;
            }
            get connected() {
                return null != this.client && this.client.readyState === WebSocket.OPEN;
            }
            readBytes(t, e = 0, r = 0) {
                this.inputByteArray.readBytes(t, e, r);
            }
            writeBytes(t, e = 0, r = 0) {
                this.outputByteArray.writeBytes(t, e, r);
            }
            writeIntLeb(t) {
                this.outputByteArray.writeIntLeb(t);
            }
            readIntLeb() {
                return this.inputByteArray.readIntLeb();
            }
            flush() {
                const t = new Uint8Array(this.outputByteArray.length);
                let e = this.outputByteArray.position;
                this.outputByteArray.position = 0;
                let r = 0;
                for (; e-- > 0; ) t[r++] = this.outputByteArray.readByte();
                this.outputByteArray.clear(), this.connected && null != this.client && this.client.send(t.buffer);
            }
            close() {
                this.internalClose();
            }
            connect(t, e) {
                return (
                    console.log("Connect to " + t + ":" + e),
                    new Promise((r, n) => {
                        (this.client = new WebSocket(t + ":" + e)),
                            (this.client.binaryType = "arraybuffer"),
                            (this.client.onmessage = this.onMessage),
                            (this.client.onclose = n),
                            (this.client.onerror = n),
                            (this.client.onopen = () => {
                                console.log("Socket connection opened"), null != this && null != this.client && ((this.client.onclose = this.onClose), (this.client.onerror = this.onError), r());
                            });
                    })
                );
            }
            internalClose() {
                console.log("Socket connection close"), null != this.client && ((this.client.onopen = null), this.client.close()), (this.client = null);
            }
            onClose() {
                console.log("Socket connection closed"), this.outputByteArray.clear(), this.onConnectionClosed.call();
            }
            onError() {
                console.log("Socket connection error"), this.outputByteArray.clear(), this.onErrorOccured.call();
            }
            onMessage(t) {
                0 === this.bytesAvailable && this.inputByteArray.clear(), this.inputByteArray.setArrayBuffer(t.data), (this.inputByteArray.position = 0), this.onDataReceived.call();
            }
        }
        n([o.boundMethod], a.prototype, "onClose", null), n([o.boundMethod], a.prototype, "onError", null), n([o.boundMethod], a.prototype, "onMessage", null), (e.Socket = a);
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.Value = e.PromiseValue = e.MultiSlotVoid = e.MultiSlot = e.Slot = e.SlotVoid = void 0);
        class n {
            constructor() {
                this.callbacks = [];
            }
            get isEmpty() {
                return 0 == this.callbacks.length;
            }
            call() {
                for (const t of this.callbacks.slice()) t();
            }
            wait() {
                return new Promise((t) => this.addOnce(t));
            }
            clear() {
                this.callbacks.length = 0;
            }
            add(t) {
                this.callbacks.indexOf(t) >= 0 || this.callbacks.push(t);
            }
            addOnce(t) {
                const e = () => {
                    t(), this.remove(e);
                };
                this.add(e);
            }
            remove(t) {
                return this.callbacks.indexOf(t) >= 0 && (this.callbacks.splice(this.callbacks.indexOf(t), 1), !0);
            }
        }
        e.SlotVoid = n;
        class o {
            constructor() {
                this.callbacks = [];
            }
            get isEmpty() {
                return 0 == this.callbacks.length;
            }
            call(t) {
                for (const e of this.callbacks.slice()) e(t);
            }
            clear() {
                this.callbacks.length = 0;
            }
            has(t) {
                return this.callbacks.indexOf(t) >= 0;
            }
            add(t) {
                this.has(t) || this.callbacks.push(t);
            }
            wait(t = () => !0) {
                return new Promise((e) => {
                    const r = (n) => {
                        t(n) && (this.remove(r), e(n));
                    };
                    this.add(r);
                });
            }
            addFilter(t, e) {
                this.callbacks.push((r) => {
                    t(r) && e(r);
                });
            }
            addOnce(t) {
                const e = (r) => {
                    t(r), this.remove(e);
                };
                this.add(e);
            }
            remove(t) {
                const e = this.callbacks.indexOf(t);
                -1 !== e && this.callbacks.splice(e, 1);
            }
        }
        e.Slot = o;
        e.MultiSlot = class {
            constructor() {
                this.slots = new Map();
            }
            call(t, e) {
                const r = this.slots.get(t);
                null != r && r.call(e);
            }
            add(t, e) {
                let r = this.slots.get(t);
                null == r && ((r = new o()), this.slots.set(t, r)), r.add(e);
            }
            remove(t, e) {
                const r = this.slots.get(t);
                null != r && r.remove(e);
            }
        };
        e.MultiSlotVoid = class {
            constructor() {
                this.slots = new Map();
            }
            call(t) {
                const e = this.slots.get(t);
                null != e && e.call();
            }
            callAll() {
                this.slots.forEach((t) => t.call());
            }
            isEmpty(t) {
                const e = this.slots.get(t);
                return null == e || e.isEmpty;
            }
            add(t, e) {
                let r = this.slots.get(t);
                null == r && ((r = new n()), this.slots.set(t, r)), r.add(e);
            }
            remove(t, e) {
                const r = this.slots.get(t);
                null != r && r.remove(e);
            }
        };
        e.PromiseValue = class {
            constructor() {
                (this.resolved = null),
                    (this.promise = new Promise((t) => {
                        this.resolved = (e) => t(e);
                    }));
            }
            resolve(t) {
                null != this.resolved && this.resolved(t);
            }
            then(t, e) {
                return this.promise.then(t, e);
            }
        };
        e.Value = class {
            constructor(t, e = !1) {
                (this._value = t), (this.forceMode = e), (this.callbacks = []);
            }
            get value() {
                return this._value;
            }
            set value(t) {
                (this._value !== t || this.forceMode) && ((this._value = t), this.dispatch());
            }
            dispatch() {
                this.call(this._value);
            }
            wait(t) {
                return new Promise((e) => {
                    const r = (n) => {
                        null != n && t(n) && (this.remove(r), e(n));
                    };
                    this.add(r);
                });
            }
            add(t, e = !1) {
                -1 == this.callbacks.indexOf(t) && (this.callbacks.push(t), e && t(this.value));
            }
            subscribe(t) {
                return this.add(t), () => this.remove(t);
            }
            addOnce(t) {
                let e = !1;
                const r = (n) => {
                    e || ((e = !0), t(n)), this.remove(r);
                };
                this.add(r);
            }
            remove(t) {
                const e = this.callbacks.indexOf(t);
                e < 0 || this.callbacks.splice(e, 1);
            }
            call(t) {
                for (const e of this.callbacks.slice()) e(t);
            }
        };
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.PacketClient = void 0);
        const n = r(0);
        class o extends n.ByteArray {
            constructor(t, e) {
                if ((super(), (this._type = 0), t <= 0 || t >= o.FORMATS.length)) throw new Error("Unknown client packet type " + t);
                (this._type = t), this.writeIntLeb(t), this.writeByte(e);
            }
            get type() {
                return this._type;
            }
            load(t) {
                const e = o.FORMATS[this._type];
                let r = !1,
                    n = 0,
                    i = 0,
                    s = t;
                for (let o = 0; o < e.length; o++) {
                    const a = e.charAt(o);
                    if ("," === a) {
                        if (r || 0 !== n) throw new Error("Bad signature for client packet " + this._type);
                        r = !0;
                        continue;
                    }
                    if ("]" === a) {
                        if (0 === n) throw new Error("Bad signature for client packet " + this._type);
                        if ((i--, 0 !== i)) {
                            o = n - 1;
                            continue;
                        }
                        (n = 0), (s = t);
                        continue;
                    }
                    if ("[" === a) {
                        if (0 !== n) throw new Error("Bad signature for server packet " + this._type);
                        const a = this.getGroupLast(e, o);
                        if (((s = t.shift()), r && null == s)) break;
                        const u = a - o - 1;
                        if (s.length % u != 0) throw new Error("Group incomplete for client packet " + this._type);
                        if (((i = s.length / u), this.writeIntLeb(i), 0 !== i)) {
                            n = o + 1;
                            continue;
                        }
                        (o = a), (s = t);
                        continue;
                    }
                    if (0 === s.length) {
                        if (r && 0 === n) break;
                        throw new Error("No data for client packet " + this._type);
                    }
                    const u = s.shift();
                    switch (a) {
                        case "S":
                            this.writeUTFLeb(u);
                            break;
                        case "I":
                            this.writeIntLeb(u);
                            break;
                        case "B":
                            this.writeByte(u);
                    }
                }
                if (t.length) throw new Error("Data " + t.length + " left in client packet " + this._type);
            }
            getGroupLast(t, e) {
                for (let r = e + 1; r < t.length; r++)
                    if ("]" === t.charAt(r)) {
                        if (r === e + 1) break;
                        return r;
                    }
                throw new Error("Bad signature for client packet " + this._type);
            }
        }
        (e.PacketClient = o),
            (o.FORMATS = [
                "",
                "I",
                "IS",
                "S",
                "IBBS,BSIIBSBSBS",
                "",
                "IIIIBI,B",
                "",
                "[I]I,I",
                "[I]BI,I",
                "IBSS",
                "I",
                "SSSB",
                "S",
                "S",
                "I",
                "B",
                "I",
                "I",
                "B",
                "I,II",
                "IB",
                "I",
                "BB",
                "[I]",
                "BB",
                "B,B",
                "",
                "I",
                "B",
                "I",
                "I",
                "B",
                "",
                "",
                "",
                "",
                "BS,II",
                "IS",
                "I",
                "I",
                "I",
                "I",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "[I]",
                "",
                "",
                "",
                "I",
                "I",
                "",
                "",
                "IBS",
                "IB",
                "I",
                "[I]",
                "[I]",
                "I",
                "BB",
                "I",
                "",
                "I",
                "",
                "",
                "B",
                "I",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "I",
                "",
                "",
                "",
                "BIS",
                "II",
                "II",
                "",
                "III",
                "",
                "",
                "[IBIIB]",
                "[I]B",
                "",
                "",
                "",
                "I",
                "",
                "BIIB",
                "",
                "",
                "",
                "I",
                "",
                "",
                "III",
                "I",
                "B",
                "B",
                "",
                "BIIB",
                "",
                "",
                "S",
                "",
                "",
                "",
                "",
                "",
                "B",
                "",
                "BII",
                "",
                "",
                "",
                "I",
                "S",
                "",
                "",
                "",
                "I",
                "BB",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "B",
                "[SS]",
                "",
                "",
                "IB",
                "",
                "",
                "",
                "B",
                "",
                "",
                "",
                "[S]",
                "SSSSII",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "B",
                "",
                "",
                "BI,BB",
                "BII,B",
                "S",
                "B",
                "",
                "",
                "",
                "",
                "",
                "BBI",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "BIBBBIB",
                "",
                "I",
                "I",
                "BS,B",
                "B,I",
                "IB",
                "IBB",
                "",
                "BB",
                "BB",
                "",
                "I",
                "IB",
                "BB",
                "",
                "S",
                "ISS",
                "",
                "",
                "I",
                "",
                "",
                "B",
                "",
                "I",
                "",
                "",
                "",
                "BB",
                "B",
                "IBSB",
                "",
                "B",
                "SIS",
                "",
                "",
                "SB",
                "II",
                "III",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "I[B],S",
                "",
                "",
                "",
                "IS",
                "",
                "",
                "I",
                "I",
                "I",
                "S",
                "BIS",
                "I",
                "B",
                "S",
                "I",
                "BII",
                "BS",
                "IS",
                "I",
                "",
                "",
                "ISS",
                "B",
                "B",
                "I",
            ]);
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.PacketServer = void 0);
        const n = r(18);
        class o extends Array {
            constructor(t, e, r) {
                super(), (this._type = 0), (this._id = 0), (this.bytesLength = 0), (this.bytesLength = r.bytesAvailable), Object.setPrototypeOf(this, o.prototype), (this._type = t), (this._id = e);
                const i = o.FORMATS[this._type];
                try {
                    n.PacketParser.readData(r, i, this, this._type);
                } catch (t) {
                    console.error(`Failed to read server packet ${this._type} "${i}": ` + t);
                }
            }
            get type() {
                return this._type;
            }
            get id() {
                return this._id;
            }
        }
        (e.PacketServer = o),
            (o.FORMATS = [
                "",
                "",
                "IIIIIII[BII]IIBBIBIII[BBBII]",
                "S",
                "B,IIII[B]IIIISBBIIBS",
                "AI,I",
                "AI,I",
                "I,B",
                "",
                "BIIIIII",
                "B,B",
                "[IFI]",
                "[BIIII]",
                "[II]",
                "[III]",
                "",
                "IIII",
                "BB",
                "BBI[II]",
                "IB",
                "I,B",
                "I",
                "B,I",
                "II",
                "B",
                "III[I][I]",
                "IB",
                "I",
                "I",
                "II,II",
                "IB",
                "II",
                "II",
                "IIBB",
                "IB",
                "",
                "",
                "BIIS,II",
                "ISBIB",
                "[IIIB]",
                "I",
                "[II]",
                "BBB",
                "[II][II][II]",
                "IB",
                "[IB]",
                "[I]",
                "[II]",
                "BI[II]",
                "BIII",
                "[I]I",
                "I,IBSIBII",
                "II",
                "IB",
                "[I]",
                "[B]",
                "BIIIB",
                "[BI]",
                "[B]",
                "",
                "",
                "",
                "",
                "",
                "",
                "[II]",
                "BI",
                "BII,I",
                "BII",
                "BI",
                "I",
                "[IIBSBI]",
                "[IB]",
                "[IIIBIIIIIIB]",
                "[I[BIII]]",
                "I[B]",
                "IBB",
                "[II]",
                "",
                "",
                "BI,[II]",
                "BIB",
                "IB",
                "I",
                "I",
                "B",
                "II,I",
                "III[IB]",
                "B",
                "I",
                "BI",
                "BI",
                "",
                "II",
                "II",
                "I[II]",
                "II",
                "I",
                "AI",
                "I[IBI]",
                "IB",
                "IB",
                "",
                "[IIIBIIBB]",
                "BIIIII",
                "[IBIIB]",
                "II",
                "[I]B",
                "",
                "I",
                "I",
                "IBIIB",
                "B",
                "IIIB[IBII]",
                "",
                "",
                "BI[II]",
                "[IIII]",
                "BIIII",
                "",
                "B[IISI[I[I]]]",
                "BB,B",
                "IBB",
                "IIB",
                "I[II]",
                "B",
                "[IB]",
                "[BI]",
                "II[I]",
                "B",
                "I",
                "I",
                "[BI]",
                "[IB]",
                "BIBIII",
                "[II]",
                "II",
                "[BB]",
                "",
                "",
                "",
                "",
                "",
                "IS",
                "BIS",
                "[B]",
                "S",
                "",
                "",
                "",
                "[IIBBS[B]IIIIIB]",
                "BIIII",
                "BIIBBI",
                "IBII",
                "IBB",
                "IIB",
                "[BI]",
                "III",
                "I[II]",
                "",
                "[I]",
                "I",
                "IIBBI",
                "[III]",
                "I",
                "IIB",
                "BB,BII",
                "",
                "S",
                "",
                "II",
                "I",
                "BI",
                "[BI]",
                "BB",
                "[I[II]]",
                "",
                "",
                "",
                "",
                "",
                "II",
                "ISSIIBII",
                "BBI[II]",
                "I",
                "IBI",
                "III[I]",
                "I",
                "IB",
                "I",
                "II",
                "BI",
                "[IB]",
                "[IB]",
                "[I]",
                "B",
                "II",
                "BB",
                "BI[II]",
                "ISS",
                "I",
                "ISB",
                "IS",
                "IS",
                "SS",
                "B",
                "III",
                "[BB]",
                "IB",
                "",
                "",
                "",
                "IB",
                "IBBIB",
                "I,B",
                "IB",
                "I[II]",
                "IBB",
                "[II]",
                "I",
                "I",
                "B",
                "",
                "[BI]",
                "IBIB[BB]",
                "[BI]",
                "BB",
                "IBI",
                "I[II]",
                "I",
                "[IIB]",
                "IBI",
                "I",
                "I[I]",
                "[I]",
                "I",
                "BB[BI]",
                "[BI]",
                "[IBS]",
                "B",
                "",
                "IIB",
                "BBB[B]",
                "IB",
                "",
                "",
                "[BI]",
                "III",
                "BBBB",
                "[II]",
                "I",
                "",
                "BBI",
                "I",
                "[BB[I]]",
                "IB",
                "I[IBSIBII]",
                "I[II]",
                "",
                "B[I]",
                "",
                "",
                "[I]",
                "IB",
                "I[BBBII]",
                "[BB]",
                "ISBBI",
                "II",
                "",
                "II",
                "II",
                ",I",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "",
                "II[BI]",
                "IB[BB]I",
                "BB",
                "B,BI",
                ",I",
                "I",
                "[BI]",
                "BI",
                "[BI]",
                ",I",
                "[II]",
                "[BI]",
                "I",
                "[II]",
                "[IIII]",
                "[SSSSII]",
                "ISIIS",
                "[S][S]",
                "[ISIS]",
                "IB",
                "I",
                "BI",
                "I",
                "",
                "BIS,II",
                "[S]",
                "[II]",
                "[III]",
                "II",
                "[BII]",
                "BI[IIII]",
                "[BBI]",
                "[BB]",
                "I",
                "II",
                "",
                "[SI]",
                "I[IBI]",
                "B[B]",
                "BI[II],I",
                "[II]",
                "BI[II][BI[I]],BIB",
                "IIB",
                "I",
                "BI[II],I",
            ]);
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }), (e.PacketParser = void 0);
        const n = r(0);
        class o {
            static readData(t, e, r, i) {
                let s = !1;
                const a = [r];
                for (let u = 0; u < e.length; u++) {
                    const l = e.charAt(u),
                        I = a[a.length - 1];
                    if ("," !== l)
                        if ("]" !== l) {
                            if (0 === t.bytesAvailable) {
                                if (s && I === r) break;
                                throw new Error("No data for server packet " + i);
                            }
                            if ("[" !== l)
                                switch (l) {
                                    case "A": {
                                        const e = t.readIntLeb(),
                                            r = new n.ByteArray();
                                        0 !== e && t.readBytes(r, 0, e), I.push(r);
                                        break;
                                    }
                                    case "S":
                                        I.push(t.readLebUTF());
                                        break;
                                    case "F":
                                        I.push(t.readFloat());
                                        break;
                                    case "I":
                                        I.push(t.readIntLeb());
                                        break;
                                    case "B":
                                        I.push(t.readUnsignedByte());
                                }
                            else {
                                const r = [];
                                if (((r.group_total = r.group_length = t.readIntLeb()), I.push(r), 0 !== r.group_length)) {
                                    (r.group_pos = u), a.push(r);
                                    continue;
                                }
                                u = o.getGroupLast(e, u);
                            }
                        } else {
                            if (I === r) throw new Error("Bad signature 3 for server packet");
                            if ((I.group_length--, 0 !== I.group_length)) {
                                u = I.group_pos;
                                continue;
                            }
                            o.splitGroup(I), delete I.group_length, delete I.group_pos, delete I.group_total, a.pop();
                        }
                    else {
                        if (s || I !== r) throw new Error("Bad signature 2 for server packet");
                        s = !0;
                    }
                }
                if (0 !== t.bytesAvailable) throw new Error("Data left in server packet " + i + " Avalible: " + t.bytesAvailable);
            }
            static splitGroup(t) {
                if (t.group_total === t.length || !this.split) return;
                const e = t.length / t.group_total,
                    r = [],
                    n = t.group_total;
                for (let o = 0; o < n; o++) r.push(t.slice(o * e, (o + 1) * e));
                (t.length = 0), t.push(...r);
            }
            static getGroupLast(t, e) {
                let r = 1;
                for (let n = e + 1; n < t.length; n++) {
                    switch (t.charAt(n)) {
                        case "]":
                            r--;
                            break;
                        case "[":
                            r++;
                    }
                    if (0 === r) {
                        if (n === e + 1) break;
                        return n;
                    }
                }
                throw new Error("Bad signature 1 for server packet");
            }
        }
        (e.PacketParser = o), (o.split = !0);
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }),
            (e.ServerPacketType = void 0),
            (function (t) {
                (t[(t.NULL = 0)] = "NULL"),
                    (t[(t.HELLO = 1)] = "HELLO"),
                    (t[(t.ADMIN_INFO = 2)] = "ADMIN_INFO"),
                    (t[(t.ADMIN_MESSAGE = 3)] = "ADMIN_MESSAGE"),
                    (t[(t.LOGIN = 4)] = "LOGIN"),
                    (t[(t.INFO = 5)] = "INFO"),
                    (t[(t.INFO_NET = 6)] = "INFO_NET"),
                    (t[(t.BALANCE = 7)] = "BALANCE"),
                    (t[(t.BUY = 9)] = "BUY"),
                    (t[(t.CONTEST_ITEMS = 10)] = "CONTEST_ITEMS"),
                    (t[(t.ADMIN_ORDERS_INFO = 11)] = "ADMIN_ORDERS_INFO"),
                    (t[(t.EVENTS = 12)] = "EVENTS"),
                    (t[(t.REWARDS = 13)] = "REWARDS"),
                    (t[(t.ADMIRERS = 14)] = "ADMIRERS"),
                    (t[(t.GIFT = 16)] = "GIFT"),
                    (t[(t.BONUS = 17)] = "BONUS"),
                    (t[(t.LEADERBOARDS = 18)] = "LEADERBOARDS"),
                    (t[(t.VIP = 20)] = "VIP"),
                    (t[(t.ROOM_INVITE = 21)] = "ROOM_INVITE"),
                    (t[(t.MOVE = 22)] = "MOVE"),
                    (t[(t.BOTTLE_PLAY_DENIED = 24)] = "BOTTLE_PLAY_DENIED"),
                    (t[(t.BOTTLE_ROOM = 25)] = "BOTTLE_ROOM"),
                    (t[(t.BOTTLE_JOIN = 26)] = "BOTTLE_JOIN"),
                    (t[(t.BOTTLE_LEAVE = 27)] = "BOTTLE_LEAVE"),
                    (t[(t.BOTTLE_LEADER = 28)] = "BOTTLE_LEADER"),
                    (t[(t.BOTTLE_ROLL = 29)] = "BOTTLE_ROLL"),
                    (t[(t.BOTTLE_KISS = 30)] = "BOTTLE_KISS"),
                    (t[(t.BOTTLE_ENTER = 35)] = "BOTTLE_ENTER"),
                    (t[(t.CHAT_MESSAGE = 37)] = "CHAT_MESSAGE"),
                    (t[(t.CHAT_WHISPER = 38)] = "CHAT_WHISPER"),
                    (t[(t.HISTORY_CONTACTS = 39)] = "HISTORY_CONTACTS"),
                    (t[(t.TITLE_POINTS = 40)] = "TITLE_POINTS"),
                    (t[(t.IGNORE_LIST = 41)] = "IGNORE_LIST"),
                    (t[(t.CHAMPION_CHEST = 42)] = "CHAMPION_CHEST"),
                    (t[(t.BEST = 43)] = "BEST"),
                    (t[(t.MONEY_BOX = 44)] = "MONEY_BOX"),
                    (t[(t.FRIENDS = 45)] = "FRIENDS"),
                    (t[(t.TOP = 46)] = "TOP"),
                    (t[(t.LEAGUES_RATING = 47)] = "LEAGUES_RATING"),
                    (t[(t.LEAGUE_GROUP = 48)] = "LEAGUE_GROUP"),
                    (t[(t.LEAGUE_INFO = 49)] = "LEAGUE_INFO"),
                    (t[(t.SEARCH = 50)] = "SEARCH"),
                    (t[(t.LAST_MESSAGE = 51)] = "LAST_MESSAGE"),
                    (t[(t.LEAGUES_TIMEOUTS = 52)] = "LEAGUES_TIMEOUTS"),
                    (t[(t.BOTTLE_LIVE_ROOMS = 54)] = "BOTTLE_LIVE_ROOMS"),
                    (t[(t.ADVENT_CALENDAR = 55)] = "ADVENT_CALENDAR"),
                    (t[(t.MESSAGE_REACTION = 56)] = "MESSAGE_REACTION"),
                    (t[(t.GIFTS_FOR_ACTIONS_STATS = 65)] = "GIFTS_FOR_ACTIONS_STATS"),
                    (t[(t.RATING_SIZE = 66)] = "RATING_SIZE"),
                    (t[(t.WEDDING_PROPOSAL_ANSWER = 67)] = "WEDDING_PROPOSAL_ANSWER"),
                    (t[(t.WEDDING_PROPOSAL_CANCEL = 68)] = "WEDDING_PROPOSAL_CANCEL"),
                    (t[(t.WEDDING_PROPOSAL_MAKE = 69)] = "WEDDING_PROPOSAL_MAKE"),
                    (t[(t.WEDDING_PROPOSAL_REFUSE = 70)] = "WEDDING_PROPOSAL_REFUSE"),
                    (t[(t.WEDDING_PROPOSAL_INFO = 71)] = "WEDDING_PROPOSAL_INFO"),
                    (t[(t.WEDDING_ADMISSIONS = 72)] = "WEDDING_ADMISSIONS"),
                    (t[(t.WEDDING_INFO = 73)] = "WEDDING_INFO"),
                    (t[(t.WEDDING_ITEMS = 75)] = "WEDDING_ITEMS"),
                    (t[(t.WEDDING_SETTLED = 76)] = "WEDDING_SETTLED"),
                    (t[(t.WEDDING_GUESTS = 77)] = "WEDDING_GUESTS"),
                    (t[(t.WEDDING_CONTEST = 80)] = "WEDDING_CONTEST"),
                    (t[(t.WEDDING_JOIN = 81)] = "WEDDING_JOIN"),
                    (t[(t.WEDDING_KISS = 82)] = "WEDDING_KISS"),
                    (t[(t.WEDDING_LEADER = 83)] = "WEDDING_LEADER"),
                    (t[(t.WEDDING_LEAVE = 84)] = "WEDDING_LEAVE"),
                    (t[(t.WEDDING_PLAY_DENIED = 85)] = "WEDDING_PLAY_DENIED"),
                    (t[(t.WEDDING_ROLL = 86)] = "WEDDING_ROLL"),
                    (t[(t.WEDDING_ROOM = 87)] = "WEDDING_ROOM"),
                    (t[(t.WEDDING_STATUS = 88)] = "WEDDING_STATUS"),
                    (t[(t.WEDDING_VOW = 89)] = "WEDDING_VOW"),
                    (t[(t.WEDDING_GARTER = 90)] = "WEDDING_GARTER"),
                    (t[(t.WEDDING_BOUQUET = 91)] = "WEDDING_BOUQUET"),
                    (t[(t.WEDDING_HAPPY = 93)] = "WEDDING_HAPPY"),
                    (t[(t.WEDDING_DIVORCE = 94)] = "WEDDING_DIVORCE"),
                    (t[(t.WEDDING_RATING_HAPPY = 95)] = "WEDDING_RATING_HAPPY"),
                    (t[(t.WEDDING_CANCEL = 96)] = "WEDDING_CANCEL"),
                    (t[(t.WEDDING_START_TIME = 97)] = "WEDDING_START_TIME"),
                    (t[(t.HOUSE_GIFT = 104)] = "HOUSE_GIFT"),
                    (t[(t.ACHIEVEMENT_GET = 106)] = "ACHIEVEMENT_GET"),
                    (t[(t.CURIOS = 117)] = "CURIOS"),
                    (t[(t.CURIOS_GIFT = 118)] = "CURIOS_GIFT"),
                    (t[(t.CHAT_HISTORY = 120)] = "CHAT_HISTORY"),
                    (t[(t.COLLECTIONS_ASSEMBLE = 121)] = "COLLECTIONS_ASSEMBLE"),
                    (t[(t.COLLECTIONS_AWARD = 122)] = "COLLECTIONS_AWARD"),
                    (t[(t.STATUS_GIFT_STATS = 127)] = "STATUS_GIFT_STATS"),
                    (t[(t.HALLOWEEN_DATA = 128)] = "HALLOWEEN_DATA"),
                    (t[(t.COLLECTIONS_POINTS = 130)] = "COLLECTIONS_POINTS"),
                    (t[(t.SELF_RICH_UPDATE = 131)] = "SELF_RICH_UPDATE"),
                    (t[(t.CHRISTMAS_DECORATIONS = 132)] = "CHRISTMAS_DECORATIONS"),
                    (t[(t.HOUSE_PACK_GIFT = 134)] = "HOUSE_PACK_GIFT"),
                    (t[(t.POSTING_REWARDS = 135)] = "POSTING_REWARDS"),
                    (t[(t.TRAINING = 145)] = "TRAINING"),
                    (t[(t.XSOLLA_SIGNATURE = 146)] = "XSOLLA_SIGNATURE"),
                    (t[(t.OFFERS_INFO = 163)] = "OFFERS_INFO"),
                    (t[(t.COINS_LADDER = 164)] = "COINS_LADDER"),
                    (t[(t.TIMEOUTS = 173)] = "TIMEOUTS"),
                    (t[(t.COUNTER = 175)] = "COUNTER"),
                    (t[(t.RATING = 183)] = "RATING"),
                    (t[(t.ADMIRE_SERIES = 200)] = "ADMIRE_SERIES"),
                    (t[(t.PLAYER_ROOM_TYPE = 212)] = "PLAYER_ROOM_TYPE"),
                    (t[(t.PHOTOS_INFO = 213)] = "PHOTOS_INFO"),
                    (t[(t.RATING_VIEWS = 216)] = "RATING_VIEWS"),
                    (t[(t.EMOTION = 217)] = "EMOTION"),
                    (t[(t.PLAYERS_KISSES = 218)] = "PLAYERS_KISSES"),
                    (t[(t.PLAYERS_VIEW = 230)] = "PLAYERS_VIEW"),
                    (t[(t.UPDATE_INFO = 232)] = "UPDATE_INFO"),
                    (t[(t.QUEST = 236)] = "QUEST"),
                    (t[(t.MODERATION_LIST = 238)] = "MODERATION_LIST"),
                    (t[(t.OFFERS_BALANCE = 249)] = "OFFERS_BALANCE"),
                    (t[(t.SUBSCRIPTION_DATA = 250)] = "SUBSCRIPTION_DATA"),
                    (t[(t.BIRTHDAY_NOTIFY = 253)] = "BIRTHDAY_NOTIFY"),
                    (t[(t.POPULAR_GIFTS = 254)] = "POPULAR_GIFTS"),
                    (t[(t.IGNORED = 255)] = "IGNORED"),
                    (t[(t.HISTORY_MESSAGES = 256)] = "HISTORY_MESSAGES"),
                    (t[(t.ROULETTE = 262)] = "ROULETTE"),
                    (t[(t.BANS = 264)] = "BANS"),
                    (t[(t.ADMIN_BUYINGS_INFO = 295)] = "ADMIN_BUYINGS_INFO"),
                    (t[(t.VIDEO_INFO = 296)] = "VIDEO_INFO"),
                    (t[(t.VIDEO_ROOM = 297)] = "VIDEO_ROOM"),
                    (t[(t.VIDEO_LISTS = 298)] = "VIDEO_LISTS"),
                    (t[(t.VIDEO_QUEUE = 299)] = "VIDEO_QUEUE"),
                    (t[(t.CHAT_OFFER = 300)] = "CHAT_OFFER"),
                    (t[(t.WINK = 301)] = "WINK"),
                    (t[(t.SPECTATOR_JOIN_LEAVE = 302)] = "SPECTATOR_JOIN_LEAVE"),
                    (t[(t.CAPTCHA = 304)] = "CAPTCHA"),
                    (t[(t.CHAT_GIF = 305)] = "CHAT_GIF"),
                    (t[(t.KISS_PRIORITY = 307)] = "KISS_PRIORITY"),
                    (t[(t.KICK_KICKS = 308)] = "KICK_KICKS"),
                    (t[(t.KICK_SAVE = 309)] = "KICK_SAVE"),
                    (t[(t.BALANCE_ITEMS = 310)] = "BALANCE_ITEMS"),
                    (t[(t.TITLES = 311)] = "TITLES"),
                    (t[(t.FRAMES = 314)] = "FRAMES"),
                    (t[(t.REWARD_GOT = 315)] = "REWARD_GOT"),
                    (t[(t.POPULAR_VIDEOS = 317)] = "POPULAR_VIDEOS"),
                    (t[(t.ADMIN_REWARDS_INFO = 318)] = "ADMIN_REWARDS_INFO"),
                    (t[(t.GIFT_BOXES = 319)] = "GIFT_BOXES"),
                    (t[(t.BALLOONS = 320)] = "BALLOONS"),
                    (t[(t.LIMITED_OFFERS = 321)] = "LIMITED_OFFERS"),
                    (t[(t.CONTEST_WITH_ACTIONS = 322)] = "CONTEST_WITH_ACTIONS"),
                    (t[(t.CONTEST_ACTION = 323)] = "CONTEST_ACTION"),
                    (t[(t.LEAGUE_CURRENT_POINTS = 324)] = "LEAGUE_CURRENT_POINTS"),
                    (t[(t.CONTEST_TETRIS = 325)] = "CONTEST_TETRIS"),
                    (t[(t.MAX_TYPE = 326)] = "MAX_TYPE");
            })(e.ServerPacketType || (e.ServerPacketType = {}));
    },
    function (t, e, r) {
        "use strict";
        Object.defineProperty(e, "__esModule", { value: !0 }),
            (e.ClientPacketType = void 0),
            (function (t) {
                (t[(t.ADMIN_REQUEST = 1)] = "ADMIN_REQUEST"),
                    (t[(t.ADMIN_EDIT = 2)] = "ADMIN_EDIT"),
                    (t[(t.LOGIN = 4)] = "LOGIN"),
                    (t[(t.REFILL = 5)] = "REFILL"),
                    (t[(t.BUY = 6)] = "BUY"),
                    (t[(t.REQUEST = 8)] = "REQUEST"),
                    (t[(t.REQUEST_NET = 9)] = "REQUEST_NET"),
                    (t[(t.ADMIN_ORDERS = 10)] = "ADMIN_ORDERS"),
                    (t[(t.GAME_REWARDS_GET = 11)] = "GAME_REWARDS_GET"),
                    (t[(t.INFO = 12)] = "INFO"),
                    (t[(t.PARAMS = 13)] = "PARAMS"),
                    (t[(t.CODE_CLAIM = 14)] = "CODE_CLAIM"),
                    (t[(t.INVITE = 15)] = "INVITE"),
                    (t[(t.SEX = 16)] = "SEX"),
                    (t[(t.VIEW = 17)] = "VIEW"),
                    (t[(t.BIRTHDAY = 18)] = "BIRTHDAY"),
                    (t[(t.COLOR = 19)] = "COLOR"),
                    (t[(t.COUNT = 20)] = "COUNT"),
                    (t[(t.MOVE = 21)] = "MOVE"),
                    (t[(t.INVITE_REF = 22)] = "INVITE_REF"),
                    (t[(t.LEADERBOARD_REQUEST = 23)] = "LEADERBOARD_REQUEST"),
                    (t[(t.ROOM_INVITE = 24)] = "ROOM_INVITE"),
                    (t[(t.FLAGS_SET = 25)] = "FLAGS_SET"),
                    (t[(t.BOTTLE_PLAY = 26)] = "BOTTLE_PLAY"),
                    (t[(t.BOTTLE_LEAVE = 27)] = "BOTTLE_LEAVE"),
                    (t[(t.BOTTLE_ROLL = 28)] = "BOTTLE_ROLL"),
                    (t[(t.BOTTLE_KISS = 29)] = "BOTTLE_KISS"),
                    (t[(t.BOTTLE_SAVE = 30)] = "BOTTLE_SAVE"),
                    (t[(t.BOTTLE_KICK = 31)] = "BOTTLE_KICK"),
                    (t[(t.BOTTLE_WAITER_JOIN = 33)] = "BOTTLE_WAITER_JOIN"),
                    (t[(t.BOTTLE_WAITER_LEAVE = 34)] = "BOTTLE_WAITER_LEAVE"),
                    (t[(t.CHAT_JOIN = 35)] = "CHAT_JOIN"),
                    (t[(t.CHAT_LEAVE = 36)] = "CHAT_LEAVE"),
                    (t[(t.CHAT_MESSAGE = 37)] = "CHAT_MESSAGE"),
                    (t[(t.CHAT_WHISPER = 38)] = "CHAT_WHISPER"),
                    (t[(t.HISTORY_REQUEST = 39)] = "HISTORY_REQUEST"),
                    (t[(t.HISTORY_CLEAR = 40)] = "HISTORY_CLEAR"),
                    (t[(t.UNFRIEND = 41)] = "UNFRIEND"),
                    (t[(t.LAST_MESSAGE_REQUEST = 42)] = "LAST_MESSAGE_REQUEST"),
                    (t[(t.TOP = 43)] = "TOP"),
                    (t[(t.REQUEST_LEAGUES_RATING = 44)] = "REQUEST_LEAGUES_RATING"),
                    (t[(t.REQUEST_LEAGUE = 45)] = "REQUEST_LEAGUE"),
                    (t[(t.HISTORY_CLEAR_COUNT = 46)] = "HISTORY_CLEAR_COUNT"),
                    (t[(t.BOTTLE_LEAVING_ROOMS = 54)] = "BOTTLE_LEAVING_ROOMS"),
                    (t[(t.IGNORE_ADD = 58)] = "IGNORE_ADD"),
                    (t[(t.IGNORE_REMOVE = 59)] = "IGNORE_REMOVE"),
                    (t[(t.RECEIVE_DAILY_BONUS = 61)] = "RECEIVE_DAILY_BONUS"),
                    (t[(t.WEDDING_PROPOSAL = 62)] = "WEDDING_PROPOSAL"),
                    (t[(t.WEDDING_PROPOSAL_ANSWER = 63)] = "WEDDING_PROPOSAL_ANSWER"),
                    (t[(t.WEDDING_PROPOSAL_CANCEL = 64)] = "WEDDING_PROPOSAL_CANCEL"),
                    (t[(t.WEDDING_REQUEST = 65)] = "WEDDING_REQUEST"),
                    (t[(t.WEDDING_ITEMS_GET = 67)] = "WEDDING_ITEMS_GET"),
                    (t[(t.WEDDING_ITEMS_SET = 68)] = "WEDDING_ITEMS_SET"),
                    (t[(t.WEDDING_INVITE = 69)] = "WEDDING_INVITE"),
                    (t[(t.WEDDING_GUESTS_REQUEST = 70)] = "WEDDING_GUESTS_REQUEST"),
                    (t[(t.WEDDING_PLAY = 71)] = "WEDDING_PLAY"),
                    (t[(t.WEDDING_LEAVE = 72)] = "WEDDING_LEAVE"),
                    (t[(t.WEDDING_ROLL = 73)] = "WEDDING_ROLL"),
                    (t[(t.WEDDING_KISS = 74)] = "WEDDING_KISS"),
                    (t[(t.WEDDING_START = 79)] = "WEDDING_START"),
                    (t[(t.WEDDING_VOW_ANSWER = 80)] = "WEDDING_VOW_ANSWER"),
                    (t[(t.WEDDING_GARTER_THROW = 81)] = "WEDDING_GARTER_THROW"),
                    (t[(t.WEDDING_COMPLETE = 82)] = "WEDDING_COMPLETE"),
                    (t[(t.WEDDING_BOUQUET_THROW = 83)] = "WEDDING_BOUQUET_THROW"),
                    (t[(t.WEDDING_DIVORCE = 86)] = "WEDDING_DIVORCE"),
                    (t[(t.WEDDING_RATING_HAPPY = 87)] = "WEDDING_RATING_HAPPY"),
                    (t[(t.WEDDING_CANCEL = 88)] = "WEDDING_CANCEL"),
                    (t[(t.WEDDING_START_NOW = 89)] = "WEDDING_START_NOW"),
                    (t[(t.CHAT_MESSAGE_EDIT = 91)] = "CHAT_MESSAGE_EDIT"),
                    (t[(t.CURIOS_REQUEST = 111)] = "CURIOS_REQUEST"),
                    (t[(t.CURIOS_GIFT = 112)] = "CURIOS_GIFT"),
                    (t[(t.GET_CHAT_HISTORY = 114)] = "GET_CHAT_HISTORY"),
                    (t[(t.COLLECTIONS_ASSEMBLE = 115)] = "COLLECTIONS_ASSEMBLE"),
                    (t[(t.CHAT_MESSAGE_REACTION = 117)] = "CHAT_MESSAGE_REACTION"),
                    (t[(t.VIDEO_POPULAR_GET = 124)] = "VIDEO_POPULAR_GET"),
                    (t[(t.FESTIVAL_REWARD_GET = 132)] = "FESTIVAL_REWARD_GET"),
                    (t[(t.PUSH_TOKEN = 133)] = "PUSH_TOKEN"),
                    (t[(t.GOLDEN_TICKET_REQUEST = 144)] = "GOLDEN_TICKET_REQUEST"),
                    (t[(t.SING_XS = 147)] = "SING_XS"),
                    (t[(t.COMPLAIN = 150)] = "COMPLAIN"),
                    (t[(t.VIDEO_INFO_GET = 158)] = "VIDEO_INFO_GET"),
                    (t[(t.VIDEO_INFO_ADD = 159)] = "VIDEO_INFO_ADD"),
                    (t[(t.GET_BDAY_REWARD = 163)] = "GET_BDAY_REWARD"),
                    (t[(t.PROFILE_REWARD_GET = 170)] = "PROFILE_REWARD_GET"),
                    (t[(t.CONTEST_ITEM_CREATE = 173)] = "CONTEST_ITEM_CREATE"),
                    (t[(t.SET_ACTION_TIMEOUT = 174)] = "SET_ACTION_TIMEOUT"),
                    (t[(t.CONTEST_ITEM_GIFT = 175)] = "CONTEST_ITEM_GIFT"),
                    (t[(t.UTM_SET = 176)] = "UTM_SET"),
                    (t[(t.CONTEST_TETRIS_BUILD = 177)] = "CONTEST_TETRIS_BUILD"),
                    (t[(t.RATING = 183)] = "RATING"),
                    (t[(t.SEARCH = 200)] = "SEARCH"),
                    (t[(t.BOTTLE_MOVE = 202)] = "BOTTLE_MOVE"),
                    (t[(t.REQUEST_PLAYER_ROOM_TYPE = 203)] = "REQUEST_PLAYER_ROOM_TYPE"),
                    (t[(t.PHOTOS_ADD_PHOTO = 204)] = "PHOTOS_ADD_PHOTO"),
                    (t[(t.PHOTOS_REMOVE_PHOTO = 205)] = "PHOTOS_REMOVE_PHOTO"),
                    (t[(t.PHOTOS_REQUEST = 206)] = "PHOTOS_REQUEST"),
                    (t[(t.PHOTOS_LIKE = 207)] = "PHOTOS_LIKE"),
                    (t[(t.INTEREST_SET = 209)] = "INTEREST_SET"),
                    (t[(t.PROFESSION_SET = 210)] = "PROFESSION_SET"),
                    (t[(t.RATING_VIEWS = 212)] = "RATING_VIEWS"),
                    (t[(t.EMOTION = 214)] = "EMOTION"),
                    (t[(t.EMAIL = 216)] = "EMAIL"),
                    (t[(t.ADMIN_PUSH_SEND = 217)] = "ADMIN_PUSH_SEND"),
                    (t[(t.GET_REWARD_SEARCH = 220)] = "GET_REWARD_SEARCH"),
                    (t[(t.WALL_POST = 221)] = "WALL_POST"),
                    (t[(t.SET_FRAME = 223)] = "SET_FRAME"),
                    (t[(t.BALLOON_BURST = 224)] = "BALLOON_BURST"),
                    (t[(t.QUEST_ACTION = 229)] = "QUEST_ACTION"),
                    (t[(t.MODERATION_REQUEST = 230)] = "MODERATION_REQUEST"),
                    (t[(t.MODERATION_DECISION = 231)] = "MODERATION_DECISION"),
                    (t[(t.TIMEZONE = 233)] = "TIMEZONE"),
                    (t[(t.VIDEO_PLAY = 234)] = "VIDEO_PLAY"),
                    (t[(t.VIDEO_LIKE = 237)] = "VIDEO_LIKE"),
                    (t[(t.OFFER_GIFT = 239)] = "OFFER_GIFT"),
                    (t[(t.BIRTHDAY_NOTIFY = 242)] = "BIRTHDAY_NOTIFY"),
                    (t[(t.BAN_MULTIPLE = 247)] = "BAN_MULTIPLE"),
                    (t[(t.BAN_ACCEPT = 248)] = "BAN_ACCEPT"),
                    (t[(t.ADMIN_BUYINGS_REQUEST = 251)] = "ADMIN_BUYINGS_REQUEST"),
                    (t[(t.ROULETTE_ROLL = 252)] = "ROULETTE_ROLL"),
                    (t[(t.IGNORE_REQUEST = 256)] = "IGNORE_REQUEST"),
                    (t[(t.VIDEO_CANCEL = 257)] = "VIDEO_CANCEL"),
                    (t[(t.WINK = 259)] = "WINK"),
                    (t[(t.CAPTCHA = 261)] = "CAPTCHA"),
                    (t[(t.ACTION = 262)] = "ACTION"),
                    (t[(t.VIDEO_AD_COMPLETE = 263)] = "VIDEO_AD_COMPLETE"),
                    (t[(t.GIF = 264)] = "GIF"),
                    (t[(t.GIF_WHISPER = 265)] = "GIF_WHISPER"),
                    (t[(t.GET_ADMIRE_BONUS = 266)] = "GET_ADMIRE_BONUS"),
                    (t[(t.REQUEST_OFFERS = 268)] = "REQUEST_OFFERS"),
                    (t[(t.ADMIN_REWARDS_REQUEST = 269)] = "ADMIN_REWARDS_REQUEST"),
                    (t[(t.GIFT_BOXES = 270)] = "GIFT_BOXES"),
                    (t[(t.TUTORIAL = 271)] = "TUTORIAL");
            })(e.ClientPacketType || (e.ClientPacketType = {}));
    },
]);
//# sourceMappingURL=connection_worker.js.map
