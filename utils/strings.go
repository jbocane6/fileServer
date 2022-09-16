package utils

const (
	// Argument format Info
	FormatSend      string = "Format: ./client send [path/filename] -channel [channel number]\n"
	FormatSubscribe string = "Format: subscribe {chan}"
	FormatPublish   string = "Format: send {chan} {filepath}"
	FormatReceive   string = "Format: ./client receive -channel [channel number]\n"
	// Connection info
	NewConn         string = "%v: Added new connection on channel %v\n"
	EndConn         string = "%v: A connection has terminated!\n"
	StartClientnoCh string = "%v: Starting client\n"
	StartClient     string = "%v: Starting client on channel %v\n"
	StartServer     string = "%v: Starting server...\n"
	// Transfer methods info
	ExpectedWriteName string = "Expected %d bytes for filename, write %d bytes\n"
	ExpectedFilename  string = "Expected %d bytes for filename, read %d bytes\n"
	ErrFileNameSize   string = "Error! Wrote %d bytes but length of name is %d!\n"
	ExpectedWrite     string = "Expected %d bytes for file, write %d bytes\n"
	ExpectedFile      string = "Expected %d bytes for file, read %d bytes\n"
	ExpectingFile     string = "Expecting %d bytes in file\n"
	ErrFileSize       string = "ERROR! File doesn't match expected size!"
	MsgTransfer       string = "%v: Transfer complete, expected %d bytes, wrote %d bytes\n"
	// Token file content
	GoData string = "/godata/"
)
