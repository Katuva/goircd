package utils

// IRC numeric replies as defined in RFC 1459
const (
	// Command responses
	RPL_WELCOME         = 001 // "Welcome to the Internet Relay Network <nick>!<user>@<host>"
	RPL_YOURHOST        = 002 // "Your host is <servername>, running version <ver>"
	RPL_CREATED         = 003 // "This server was created <date>"
	RPL_MYINFO          = 004 // "<servername> <version> <available user modes> <available channel modes>"
	RPL_TOPICWHOTIME    = 333 // "<channel> <who> <time>"
	RPL_BOUNCE          = 005 // "Try server <server name>, port <port number>"
	RPL_USERHOST        = 302 // Reply format used by USERHOST
	RPL_ISON            = 303 // Reply format used by ISON
	RPL_AWAY            = 301 // "<nick> :<away message>"
	RPL_UNAWAY          = 305 // ":You are no longer marked as being away"
	RPL_NOWAWAY         = 306 // ":You have been marked as being away"
	RPL_WHOISUSER       = 311 // "<nick> <user> <host> * :<real name>"
	RPL_WHOISSERVER     = 312 // "<nick> <server> :<server info>"
	RPL_WHOISOPERATOR   = 313 // "<nick> :is an IRC operator"
	RPL_WHOISIDLE       = 317 // "<nick> <integer> :seconds idle"
	RPL_ENDOFWHOIS      = 318 // "<nick> :End of WHOIS list"
	RPL_WHOISCHANNELS   = 319 // "<nick> :*( ( "@" / "+" ) <channel> " " )"
	RPL_WHOWASUSER      = 314 // "<nick> <user> <host> * :<real name>"
	RPL_ENDOFWHOWAS     = 369 // "<nick> :End of WHOWAS"
	RPL_LISTSTART       = 321 // :server 321 nick Channel :Users  Name
	RPL_LIST            = 322 // "<channel> <# visible> :<topic>"
	RPL_LISTEND         = 323 // ":End of LIST"
	RPL_UNIQOPIS        = 325 // "<channel> <nickname>"
	RPL_CHANNELMODEIS   = 324 // "<channel> <mode> <mode params>"
	RPL_NOTOPIC         = 331 // "<channel> :No topic is set"
	RPL_TOPIC           = 332 // "<channel> :<topic>"
	RPL_INVITING        = 341 // "<channel> <nick>"
	RPL_SUMMONING       = 342 // "<user> :Summoning user to IRC"
	RPL_INVITELIST      = 346 // "<channel> <invitemask>"
	RPL_ENDOFINVITELIST = 347 // "<channel> :End of channel invite list"
	RPL_EXCEPTLIST      = 348 // "<channel> <exceptionmask>"
	RPL_ENDOFEXCEPTLIST = 349 // "<channel> :End of channel exception list"
	RPL_VERSION         = 351 // "<version>.<debuglevel> <server> :<comments>"
	RPL_WHOREPLY        = 352 // "<channel> <user> <host> <server> <nick> <H|G>[*][@|+] :<hopcount> <real name>"
	RPL_ENDOFWHO        = 315 // "<name> :End of WHO list"
	RPL_NAMEREPLY       = 353 // "( "=" / "*" / "@" ) <channel> :[ "@" / "+" ] <nick> *( " " [ "@" / "+" ] <nick> )"
	RPL_ENDOFNAMES      = 366 // "<channel> :End of NAMES list"
	RPL_LINKS           = 364 // "<mask> <server> :<hopcount> <server info>"
	RPL_ENDOFLINKS      = 365 // "<mask> :End of LINKS list"
	RPL_BANLIST         = 367 // "<channel> <banmask>"
	RPL_ENDOFBANLIST    = 368 // "<channel> :End of channel ban list"
	RPL_INFO            = 371 // ":<string>"
	RPL_ENDOFINFO       = 374 // ":End of INFO list"
	RPL_MOTDSTART       = 375 // ":- <server> Message of the day - "
	RPL_MOTD            = 372 // ":- <text>"
	RPL_ENDOFMOTD       = 376 // ":End of MOTD command"
	RPL_YOUREOPER       = 381 // ":You are now an IRC operator"
	RPL_REHASHING       = 382 // "<config file> :Rehashing"
	RPL_YOURESERVICE    = 383 // "You are service <servicename>"
	RPL_TIME            = 391 // "<server> :<string showing server's local time>"
	RPL_USERSSTART      = 392 // ":UserID   Terminal  Host"
	RPL_USERS           = 393 // ":<username> <ttyline> <hostname>"
	RPL_ENDOFUSERS      = 394 // ":End of users"
	RPL_NOUSERS         = 395 // ":Nobody logged in"
	RPL_TRACELINK       = 200 // "Link <version & debug level> <destination> <next server> V<protocol version> <link uptime in seconds> <backstream sendq> <upstream sendq>"
	RPL_TRACECONNECTING = 201 // "Try. <class> <server>"
	RPL_TRACEHANDSHAKE  = 202 // "H.S. <class> <server>"
	RPL_TRACEUNKNOWN    = 203 // "???? <class> [<client IP address in dot form>]"
	RPL_TRACEOPERATOR   = 204 // "Oper <class> <nick>"
	RPL_TRACEUSER       = 205 // "User <class> <nick>"
	RPL_TRACESERVER     = 206 // "Serv <class> <int>S <int>C <server> <nick!user|*!*>@<host|server> V<protocol version>"
	RPL_TRACESERVICE    = 207 // "Service <class> <name> <type> <active type>"
	RPL_TRACENEWTYPE    = 208 // "<newtype> 0 <client name>"
	RPL_TRACECLASS      = 209 // "Class <class> <count>"
	RPL_TRACERECONNECT  = 210 // Unused
	RPL_TRACELOG        = 261 // "File <logfile> <debug level>"
	RPL_TRACEEND        = 262 // "<server name> <version & debug level> :End of TRACE"
	RPL_STATSLINKINFO   = 211 // "<linkname> <sendq> <sent messages> <sent Kbytes> <received messages> <received Kbytes> <time open>"
	RPL_STATSCOMMANDS   = 212 // "<command> <count> <byte count> <remote count>"
	RPL_ENDOFSTATS      = 219 // "<stats letter> :End of STATS report"
	RPL_STATSUPTIME     = 242 // ":Server Up %d days %d:%02d:%02d"
	RPL_STATSOLINE      = 243 // "O <hostmask> * <name>"
	RPL_UMODEIS         = 221 // "<user mode string>"
	RPL_SERVLIST        = 234 // "<name> <server> <mask> <type> <hopcount> <info>"
	RPL_SERVLISTEND     = 235 // "<mask> <type> :End of service listing"
	RPL_LUSERCLIENT     = 251 // ":There are <integer> users and <integer> services on <integer> servers"
	RPL_LUSEROP         = 252 // "<integer> :operator(s) online"
	RPL_LUSERUNKNOWN    = 253 // "<integer> :unknown connection(s)"
	RPL_LUSERCHANNELS   = 254 // "<integer> :channels formed"
	RPL_LUSERME         = 255 // ":I have <integer> clients and <integer> servers"
	RPL_ADMINME         = 256 // "<server> :Administrative info"
	RPL_ADMINLOC1       = 257 // ":<admin info>"
	RPL_ADMINLOC2       = 258 // ":<admin info>"
	RPL_ADMINEMAIL      = 259 // ":<admin info>"
	RPL_TRYAGAIN        = 263 // "<command> :Please wait a while and try again."

	// Error replies
	ERR_NOSUCHNICK        = 401 // "<nickname> :No such nick/channel"
	ERR_NOSUCHSERVER      = 402 // "<server name> :No such server"
	ERR_NOSUCHCHANNEL     = 403 // "<channel name> :No such channel"
	ERR_CANNOTSENDTOCHAN  = 404 // "<channel name> :Cannot send to channel"
	ERR_TOOMANYCHANNELS   = 405 // "<channel name> :You have joined too many channels"
	ERR_WASNOSUCHNICK     = 406 // "<nickname> :There was no such nickname"
	ERR_TOOMANYTARGETS    = 407 // "<target> :<error code> recipients. <abort message>"
	ERR_NOSUCHSERVICE     = 408 // "<service name> :No such service"
	ERR_NOORIGIN          = 409 // ":No origin specified"
	ERR_NORECIPIENT       = 411 // ":No recipient given (<command>)"
	ERR_NOTEXTTOSEND      = 412 // ":No text to send"
	ERR_NOTOPLEVEL        = 413 // "<mask> :No toplevel domain specified"
	ERR_WILDTOPLEVEL      = 414 // "<mask> :Wildcard in toplevel domain"
	ERR_BADMASK           = 415 // "<mask> :Bad Server/host mask"
	ERR_UNKNOWNCOMMAND    = 421 // "<command> :Unknown command"
	ERR_NOMOTD            = 422 // ":MOTD File is missing"
	ERR_NOADMININFO       = 423 // "<server> :No administrative info available"
	ERR_FILEERROR         = 424 // ":File error doing <file op> on <file>"
	ERR_NONICKNAMEGIVEN   = 431 // ":No nickname given"
	ERR_ERRONEUSNICKNAME  = 432 // "<nick> :Erroneous nickname"
	ERR_NICKNAMEINUSE     = 433 // "<nick> :Nickname is already in use"
	ERR_NICKCOLLISION     = 436 // "<nick> :Nickname collision KILL from <user>@<host>"
	ERR_UNAVAILRESOURCE   = 437 // "<nick/channel> :Nick/channel is temporarily unavailable"
	ERR_USERNOTINCHANNEL  = 441 // "<nick> <channel> :They aren't on that channel"
	ERR_NOTONCHANNEL      = 442 // "<channel> :You're not on that channel"
	ERR_USERONCHANNEL     = 443 // "<user> <channel> :is already on channel"
	ERR_NOLOGIN           = 444 // "<user> :User not logged in"
	ERR_SUMMONDISABLED    = 445 // ":SUMMON has been disabled"
	ERR_USERSDISABLED     = 446 // ":USERS has been disabled"
	ERR_NOTREGISTERED     = 451 // ":You have not registered"
	ERR_NEEDMOREPARAMS    = 461 // "<command> :Not enough parameters"
	ERR_ALREADYREGISTRED  = 462 // ":Unauthorized command (already registered)"
	ERR_NOPERMFORHOST     = 463 // ":Your host isn't among the privileged"
	ERR_PASSWDMISMATCH    = 464 // ":Password incorrect"
	ERR_YOUREBANNEDCREEP  = 465 // ":You are banned from this server"
	ERR_YOUWILLBEBANNED   = 466 // ""
	ERR_KEYSET            = 467 // "<channel> :Channel key already set"
	ERR_CHANNELISFULL     = 471 // "<channel> :Cannot join channel (+l)"
	ERR_UNKNOWNMODE       = 472 // "<char> :is unknown mode char to me for <channel>"
	ERR_INVITEONLYCHAN    = 473 // "<channel> :Cannot join channel (+i)"
	ERR_BANNEDFROMCHAN    = 474 // "<channel> :Cannot join channel (+b)"
	ERR_BADCHANNELKEY     = 475 // "<channel> :Cannot join channel (+k)"
	ERR_BADCHANMASK       = 476 // "<channel> :Bad Channel Mask"
	ERR_NOCHANMODES       = 477 // "<channel> :Channel doesn't support modes"
	ERR_BANLISTFULL       = 478 // "<channel> <char> :Channel list is full"
	ERR_NOPRIVILEGES      = 481 // ":Permission Denied- You're not an IRC operator"
	ERR_CHANOPRIVSNEEDED  = 482 // "<channel> :You're not channel operator"
	ERR_CANTKILLSERVER    = 483 // ":You can't kill a server!"
	ERR_RESTRICTED        = 484 // ":Your connection is restricted!"
	ERR_UNIQOPPRIVSNEEDED = 485 // ":You're not the original channel operator"
	ERR_NOOPERHOST        = 491 // ":No O-lines for your host"
	ERR_UMODEUNKNOWNFLAG  = 501 // ":Unknown MODE flag"
	ERR_USERSDONTMATCH    = 502 // ":Cannot change mode for other users"
)

// Server version and name
const (
	SERVER_NAME    = "goircd"
	SERVER_VERSION = "0.1.1"
	SERVER_CREATED = "2025-07-28"
)

// Channel modes
const (
	CHAN_MODE_PRIVATE       = 'p' // Private channel
	CHAN_MODE_SECRET        = 's' // Secret channel
	CHAN_MODE_INVITE_ONLY   = 'i' // Invite-only channel
	CHAN_MODE_TOPIC_OP_ONLY = 't' // Topic settable by channel operators only
	CHAN_MODE_NO_EXTERNAL   = 'n' // No messages from outside the channel
	CHAN_MODE_MODERATED     = 'm' // Moderated channel
	CHAN_MODE_LIMIT         = 'l' // User limit
	CHAN_MODE_KEY           = 'k' // Channel key (password)
	CHAN_MODE_OP            = 'o' // Channel operator
	CHAN_MODE_VOICE         = 'v' // Voice privilege
	CHAN_MODE_BAN           = 'b' // Ban mask
)

// User modes
const (
	USER_MODE_INVISIBLE  = 'i' // Invisible
	USER_MODE_OPERATOR   = 'o' // IRC Operator
	USER_MODE_WALLOPS    = 'w' // Receive wallops
	USER_MODE_RESTRICTED = 'r' // Restricted connection
)

// Maximum limits
const (
	MAX_CHANNELS_PER_USER   = 10
	MAX_CLIENTS_PER_CHANNEL = 100
	MAX_BANS_PER_CHANNEL    = 50
	MAX_NICK_LENGTH         = 9
	MAX_CHANNEL_NAME_LENGTH = 50
	MAX_TOPIC_LENGTH        = 307
	MAX_KICK_COMMENT_LENGTH = 307
	MAX_AWAY_MESSAGE_LENGTH = 307
)

const (
	ColourReset  = "\033[0m"
	ColourRed    = "\033[31m"
	ColourGreen  = "\033[32m"
	ColourYellow = "\033[33m"
	ColourBlue   = "\033[34m"
	ColourPurple = "\033[35m"
	ColourCyan   = "\033[36m"
	ColourWhite  = "\033[37m"
	ColourBold   = "\033[1m"

	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"
	BgWhite  = "\033[47m"
)
