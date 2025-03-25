Learn Blockchains
=================

This is a simple blockchain implementation in Go. It's a toy project to learn about blockchains.

TODO
----

- Proto block should be very generic and store the PoW data as a byte array.
  Specific block implementations should unmarshal the PoW data into a struct.
  This way we can have a generic block that can be used for any PoW algorithm.
