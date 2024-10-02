# Make the CLI a thin client of the server

previously, the cli reads from it's config file and issues commands directly against a server.

The problem with this is it's hard for it to know it's own address, which is needed when the server wants to respond.

This can be alleviated by changing the purpose of the CLI. Now it must know the address of the server running it's behalf. Aside from solving the above problem, this also provides a better guarantee that servers have something to respond to.

This means there are two fundamentally different ways of communicating with a server: As a peer sending it messages, and as an authenticated user, telling it to do things.

This can be easily acheived: both server and CLI can read from the same config, which contains an address section. The server will have read/write access, and the CLI just read. If the CLI wants to change the config file, it will tell the server to do it.

Of course all such calls can and should be authenticated. Both programs have access to the same private key.


