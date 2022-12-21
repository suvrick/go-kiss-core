import { PacketParser } from "../../../protocol/PacketParser";
import { ServerPacketType } from "../../../protocol/ServerPacketType";
import type { ByteArray } from "../../net/ByteArray";
import type { IPlayerMaskReader } from "./IPlayerMaskReader";
import { PlayerFieldsInfo } from "./PlayerFieldsInfo";
import type { PlayerInfoFields } from "./PlayerInfoFields";
import type { PlayerInfoKeys } from "./PlayerInfoKeys";

export class PlayerInfoParser
{
	private static readonly FORMATS_1: PlayerFieldsInfo = new PlayerFieldsInfo()
		.defineField("I", "nid")
		.defineField("B", "type")
		.defineField("S", "name")
		.defineField("B", "sex")
		.defineField("I", "tag")
		.defineField("I", "referrer")
		.defineField("I", "bdate")
		.defineField("SB", [ "avatar", "avatar_status" ])
		.defineField("S", "profile")
		.defineField("S", "status")
		.defineField("B", "countryId")
		.defineField("B", "online")
		.defineField("I", "admirer_id")
		.defineField("I", "admirer_price")
		.defineField("I", "admirer_time_finish")
		.defineField("I", "views")
		.defineField("B", "vip")
		.defineField("B", "color")
		.defineField("II", [ "kisses", "kisses_today" ])
		.defineField("II", [ "gifts", "gifts_today" ])
		.defineField("[III]", "lastGifts") //[source_id:I, gift_id:I, time:I]
		.defineField("B", "device")
		.defineField("I", "wedding_id")
		.defineField("[III]", "achievements")
		.defineField("[BI]", [ "collections" ])
		.defineField("B", "avatar_id")
		.defineField("B", "rights")
		.defineField("I", "register_time")
		.defineField("I", "logout_time")
		.defineField("[S][B]", [ "photos", "photos_statuses" ])
		.defineField("IIBII", [ "bridals_place", "wedlocks_place", "is_popular", "rich_place", "views_place" ])
		.defineField("B", "frame_id");

	private static readonly FORMATS_2: PlayerFieldsInfo = new PlayerFieldsInfo()
		.defineField("I", "level")
		.defineField("IB", [ "ability_expire", "ability_type" ])
		.defineField("B", "rolls_rewarded")
		.defineField("IIB", [ "subscribe_past_days", "vip_days_left", "vip_trial_used" ])
		.defineField("BIII", [ "league", "league_common_points", "league_group_id", "league_points" ])
		.defineField("I", "last_complaint_date")
		.defineField("I", "admire_reward_timestamp")
		.defineField("B", "deleted");

	public static readonly ALL: [ number, number ] = PlayerInfoParser.all();

	public static getMaskByFields<K extends PlayerInfoKeys>(fields: K[]): [ number, number ]
	{
		const result: [ number, number ] = [ 0, 0 ];
		for (const field of fields)
		{
			result[0] |= PlayerInfoParser.FORMATS_1.fields[field as string] >>> 0;
			result[1] |= PlayerInfoParser.FORMATS_2.fields[field as string] >>> 0;
		}
		return result;
	}

	public static parse(data: ByteArray, mask1: number, mask2: number): (Partial<PlayerInfoFields> & { uid: number })[]
	{
		data.position = 0;

		const parse1: [ string, IPlayerMaskReader[] ] = PlayerInfoParser.parseFormat(mask1, PlayerInfoParser.FORMATS_1);
		const parse2: [ string, IPlayerMaskReader[] ] = PlayerInfoParser.parseFormat(mask2, PlayerInfoParser.FORMATS_2);

		const format: string = `[I${parse1[0] + parse2[0]}]`;
		const fields: IPlayerMaskReader[] = parse1[1].concat(parse2[1]);

		let output: any[] = [];

		PacketParser.readData(data, format, output, ServerPacketType.INFO);

		if (output.length == 0)
			return [];

		output = output.pop();

		const result: (Partial<PlayerInfoFields> & { uid: number })[] = [];
		for (let i: number = 0; i < output.length; i++)
		{
			const rawData: any = output[i];
			let k: number = 0;
			const dataP: Partial<PlayerInfoFields> & { uid: number } = { uid: rawData[k++] };

			for (const reader of fields)
			{
				for (const field of reader.fields)
				{
					dataP[field] = rawData[k++];

					if (field == "bdate" && dataP.bdate != null && dataP.bdate >= 2147483647)
						dataP.bdate -= 0x100000000;
				}
			}

			result[i] = dataP;
		}

		return result;
	}

	private static all(): [ number, number ]
	{
		let mask1: number = 0;
		let mask2: number = 0;

		for (let i: number = 0; i < this.FORMATS_1.count; i++)
			mask1 |= 1 << i;

		for (let i = 0; i < this.FORMATS_2.count; i++)
			mask2 |= 1 << i;

		return [ mask1, mask2 ];
	}

	private static parseFormat(mask: number, formats: PlayerFieldsInfo): [ string, IPlayerMaskReader[] ]
	{
		let format: string = "";
		let fields: string[] = [];
		const info: IPlayerMaskReader[] = [];

		for (let i: number = 0; i < formats.count; i++)
		{
			const bit: number = (1 << i) >>> 0;
			if ((mask & bit) == 0)
				continue;

			info.push(formats.formats[bit]);
			format = format + formats.formats[bit].format;

			fields = fields.concat(formats.formats[bit].fields);
		}

		return [ format, info ];
	}
}