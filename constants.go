package tclient

const (
	// version constant
	VERSION				= 0.1
	// TELNET_IAC interpret as command
	TELNET_IAC			= 255
	// TELNET_DONT don't use following option
	TELNET_DONT			= 254
	// TELNET_DO use following option
	TELNET_DO			= 253
	// TELNET_WONT i will not use this option
	TELNET_WONT			= 252
	// TELNET_WILL i will use this option
	TELNET_WILL			= 251
	// TELNET_SB interpret as subnegotiation
	TELNET_SB			= 250
	// TELNET_GA you may reverse this line
	TELNET_GA			= 249
	// TELNET_EL erase current line
	TELNET_EL			= 248
	// TELNET_EC erase current character
	TELNET_EC			= 247
	// TELNET_AYT are you there
	TELNET_AYT			= 246
	// TELNET_AO abort output but let prog finish
	TELNET_AO			= 245
	// TELNET_IP interrupt process permanently
	TELNET_IP			= 244
	// TELNET_BREAK break
	TELNET_BREAK		= 243
	// TELNET_DM data mark - for connect
	TELNET_DM			= 242
	// TELNET_NOP nop
	TELNET_NOP			= 241
	// TELNET_SE subnegotiation end
	TELNET_SE			= 240
	// TELNET_EOR end of record
	TELNET_EOR			= 239
	// TELNET_ABORT abort process
	TELNET_ABORT		= 238
	// TELNET_SUSP suspend process
	TELNET_SUSP			= 237
	// TELNET_EOF end of file
	TELNET_EOF			= 236
	// TELNET_SYNC for telfunc calls
	TELNET_SYNCH		= 242

	// TELOPT_BINARY binary transmission
	TELOPT_BINARY		= 0
	// TELOPT_ECHO echo
	TELOPT_ECHO			= 1
	// TELOPT_RCP reconnetion
	TELOPT_RCP			= 2
	// TELOPT_SGA suppres go ahead
	TELOPT_SGA			= 3
	// TELOPT_NAMS approx message size negotiation
	TELOPT_NAMS			= 4
	// TELOPT_STATUS status
	TELOPT_STATUS		= 5
	// TELOPT_TM timing mark
	TELOPT_TM			= 6
	// TELOPT_RCTE remote controlled trans and echo
	TELOPT_RCTE			= 7
	// TELOPT_NAOL output line width
	TELOPT_NAOL			= 8
	// TELOPT_NAOP output page size
	TELOPT_NAOP			= 9
	// TELOPT_NAOCRD output carriage-return disposition
	TELOPT_NAOCRD		= 10
	// TELOPT_NAOHTS output horizontal tab stops
	TELOPT_NAOHTS		= 11
	// TELOPT_NAOHTD output horizontal tab disposition
	TELOPT_NAOHTD		= 12
	// TELOPT_NAOFFD  output formfeed disposition
	TELOPT_NAOFFD		= 13
	// TELOPT_NAOVTS output vertical tabstops
	TELOPT_NAOVTS		= 14
	// TELOPT_NAOVTD output vertical tab disposition
	TELOPT_NAOVTD		= 15
	// TELOPT_NAOLFD output linefeed disposition
	TELOPT_NAOLFD		= 16
	// TELOPT_XASCII extended ascii
	TELOPT_XASCII		= 17
	// TELOPT_LOGOUT logout
	TELOPT_LOGOUT		= 18
	// TELOPT_BM byte macro
	TELOPT_BM			= 19
	// TELOPT_DET data entry terminal
	TELOPT_DET			= 20
	// TELOPT_SUPDUP SUPDUP
	TELOPT_SUPDUP		= 21
	// TELOPT_SUPDUPOUTPUT SUPDUP output
	TELOPT_SUPDUPOUTPUT	= 22
	// TELOPT_SNDLOC Send location
	TELOPT_SNDLOC		= 23
	// TELOPT_TTYPE terminal type
	TELOPT_TTYPE		= 24
	// TELOPT_EOR end of record
	TELOPT_EOR			= 25
	// TELOPT_TUID TACACS user identification
	TELOPT_TUID			= 26
	// TELOPT_OUTMARK output marking
	TELOPT_OUTMARK		= 27
	// TELOPT_TTYLOC terminal location number
	TELOPT_TTYLOC		= 28
	// TELOPT_NAWS - Negotiate About Window Size
	TELOPT_NAWS			= 31
	// TELOPT_FLOWCTRL - remote flow control
	TELOPT_FLOWCTRL		= 33

	// TELOPT_SB_SEND SEND subneg
	TELOPT_SB_SEND		= 1
	// TELOPT_SB_IS IS subneg
	TELOPT_SB_IS		= 0
	// TELOPT_SB_NEV_ENVIRON
	TELOPT_SB_NEV_ENVIRON = 39
)
