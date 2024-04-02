# Test Suite

The following scenario is called Black Dust, Green Sunset, Ice Cherry, Patient Haze. It tests the following situation:

1. Ice Cherry enters. It is the first node, and all alone.
2. Green Sunset spins up. It says hi to Ice Cherry.
3. Ice Cherry says hi back.
4. Black Dust enters. It says hi to Green Sunset.
5. Green Sunset says hi back.
6. Green Sunset introduces Black Dust to Ice Cherry.
7. Patient Haze enters, saying hi to the three other nodes.
8. The nodes say hi back.
9. After some chatter, all is calm, and all nodes have three friends (each other).

To assist in debugging and observability, every node is configured to write it's log messages to a central place. let's say redis running on localhost.

