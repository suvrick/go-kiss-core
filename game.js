var format = PacketClient.FORMATS[this._type];
var optional = false;
var groupPos = 0;
var groupLength = 0;
var container = rest;
for (var i = 0; i < format.length; i++) {
    var symbol = format.charAt(i);
    if (symbol === ",") {
        if (optional || groupPos !== 0)
            throw new Error("Bad signature for client packet " + this._type);
        optional = true;
        continue;
    }
    if (symbol === "]") {
        if (groupPos === 0)
            throw new Error("Bad signature for client packet " + this._type);
        groupLength--;
        if (groupLength !== 0) {
            i = groupPos - 1;
            continue;
        }
        groupPos = 0;
        container = rest;
        continue;
    }
    if (symbol === "[") {
        if (groupPos !== 0)
            throw new Error("Bad signature for server packet " + this._type);
        var last = this.getGroupLast(format, i);
        container = rest.shift();
        if (optional && container == null)
            break;
        var count = last - i - 1;
        if (container.length % count !== 0)
            throw new Error("Group incomplete for client packet " + this._type);
        groupLength = container.length / count;
        this.writeIntLeb(groupLength);
        if (groupLength !== 0) {
            groupPos = i + 1;
            continue;
        }
        i = last;
        container = rest;
        continue;
    }
    if (container.length === 0) {
        if (optional && groupPos === 0)
            break;
        throw new Error("No data for client packet " + this._type);
    }
    var value = container.shift();
    switch (symbol) {
        case "S":
            this.writeUTFLeb(value);
            break;
        case "I":
            this.writeIntLeb(value);
            break;
        case "B":
            this.writeByte(value);
            break;
    }
}
if (rest.length)
    throw new Error("Data " + rest.length + " left in client packet " + this._type);
};
PacketClient.prototype.getGroupLast = function (format, first) {
for (var last = first + 1; last < format.length; last++) {
    if (format.charAt(last) !== "]")
        continue;
    if (last === first + 1)
        break;
    return last;
}
throw new Error("Bad signature for client packet " + this._type);
};