minimalistic http proxy

usage:

1. configure config.js where
	"port" is the port the proxy listens on
	"proxy" is another http proxy mumbojumbo
		can route traffic through in the form host:port
	"allowed" is a list of ip's that are allowed to connect
	
2. start mumbojumbo

mumbojumbo writes every request it receives into the file log
