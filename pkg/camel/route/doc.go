/*
Package route provides a set of primitives and tools to build camel route definitions.
Route is a set of processing steps that are applied to a message as it travels from a source to a destination.
A route typically consists of a series of processing steps that are connected in a linear sequence.

# errmatcher

	Provides error matcher definitions (is, any, contains, regex). Used by try step to describe catch clause.

# expr

	Provides expression definitions (simple, constant, func). Used by different steps choice, setBody, loop, etc.

# step

	Provides route step definitions:
		-	choice
		-	convertbody
		-	convertheader
		-	convertproperty
		-	delay
		-	fn
		-	log
		-	loop
		-	marshal
		-	multicast
		-	pipeline
		-	removeheader
		-	removeproperty
		-	setbody
		-	seterror
		-	setheader
		-	setproperty
		-	to
		-	try
		-	unmarshal

# Builder

	Represents a route builder.
*/
package route
