package socket

// pack, err := Marshal(payload)

// if err != nil {
// 	if socket.errorHandle != nil {
// 		socket.errorHandle(ErrMarshalClientPacket)
// 	}
// 	return
// }

// data := make([]byte, 0)

// data = AppendInt(data, int64(socket.msgID)) // message ID

// data = AppendUint(data, uint64(packType)) // packet type

// data = AppendUint(data, uint64(6)) //device

// data = append(data, pack...)

// data_len := make([]byte, 0)

// data_len = AppendInt(data_len, int64(len(data))) // packet len

// data_len = append(data_len, data...)

// socket.Logger.Printf("[SendPacket] >> %x\n", pack)

//socket.Logger.Printf("[SendPacket2] >> %x\n", data_len)

// ReadInt(reader, 32)
// 			ReadInt(reader, 32)
// 			t, _ := ReadUint(reader, 16)

// 			packetType := PacketServerType(t)

// 			//socket.Logger.Printf("server type: %d\n", packetType)

// 			switch packetType {
// 			case LOGIN_SERVER:
// 				login := PSLogin{}
// 				if err := Unmarshal(reader, &login); err != nil {
// 					if socket.errorHandle != nil {
// 						socket.errorHandle(err)
// 					}
// 				} else {
// 					socket.readHandle(packetType, login)
// 					socket.SendPacket(GET_BONUS, PCGetBonus{})
// 				}
// 			case BALANCE:
// 				balance := PSBalance{}
// 				if err := Unmarshal(reader, &balance); err != nil {
// 					if socket.errorHandle != nil {
// 						socket.errorHandle(err)
// 					}
// 				} else {
// 					socket.readHandle(packetType, balance)
// 				}
// 			case REWARDS:
// 				rewards := PSRewards{}
// 				if err := Unmarshal(reader, &rewards); err != nil {
// 					if socket.errorHandle != nil {
// 						socket.errorHandle(err)
// 					}
// 				} else {
// 					socket.readHandle(packetType, rewards)
// 				}
// 			case BONUS:
// 				bonus := PSBonus{}
// 				if err := Unmarshal(reader, &bonus); err != nil {
// 					if socket.errorHandle != nil {
// 						socket.errorHandle(err)
// 					}
// 				} else {
// 					socket.readHandle(packetType, bonus)
// 				}
// 			case BALANCE_ITEMS:
// 				balanceItems := PSBalanceItems{}
// 				if err := Unmarshal(reader, &balanceItems); err != nil {
// 					if socket.errorHandle != nil {
// 						socket.errorHandle(err)
// 					}
// 				} else {
// 					socket.readHandle(packetType, balanceItems)
// 				}
// 			}
