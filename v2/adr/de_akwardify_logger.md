# De-akwardify logger

The way logger is handled is not very elegant. I have to pass in an io.Writer to methods like NewPrincipal().

Use the option pattern instead.

So the default logger will be a nil logger that does nothing.


It will be up to the caller of NewPrincipal() to create their own logger and do something like

p := NewPrincipal( WithLogger(myLogger) )

or:

p.With( WithLogger(myLogger) )



