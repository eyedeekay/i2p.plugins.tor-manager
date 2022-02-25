A Slightly Different Take on Resistance to Blocking
===================================================

I2P has some fundamentally collaborative aspects at it's core. Routing is done using a selection
of "collaborating" peers to create a tunnel, for instance. Another way I2P is fundamentally
collaborative is in it's approach to distributing I2P updates to all I2P users.

During the update cycle, all I2P users participate in a Bittorrent swarm and assist eachother
in obtaining the update. This allows updates to be distributed more quickly, and without the
need for every downloader to use a single source, making the update essentially impossible to
tamper with or block, because as soon as the update is started, there is no single point of
failure, in fact, there is extreme redundancy. For a few days, I2P update torrents are
extraordinarily well-seeded, which is long enough to distribute the updates.

The update design in I2P can inform a plan to make Tor Browser more redundant, resilient, and
available by consistently placing the Tor Browser on the I2P torrent network. Still, we don't
necessarily want all the features of the I2P update mechanism in our tool, and there are some
features which we may want to retain or enhance.

 - Programmatically trigger an UpdatePostProcessor - Discard
 - Torrent based downloading with I2PSnark - Keep but enhance(Add BiglyBT support)
