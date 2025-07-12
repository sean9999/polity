# peer is knowledge is fact

the whole "fact base" thing was too complex. Here's the new idea:

The data structure that holds a Principal's list of peers, and the one representing the knowledge-base will be merged together.

Rather than a `map[string]Peer`, we're gonna have a `map[Peer]Record` where `Record` is our knowledge-base.

It is a struct that will grow over time, adopting new fields.

Therefore, the Option pattern will be useful.