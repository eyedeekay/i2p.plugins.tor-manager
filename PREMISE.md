Rationale for Non-Standard modifications to Tor Browser for I2P Browsing and Local Administration
=================================================================================================

This system prefers the Tor Browser for browsing the I2P network because any sane assessment of the
situation indicates that Tor Browser is the best choice for browsing HTTP/S services(web sites) in
a secure and private way at this time. However, we must make non-recommended changes to the Tor Browser
to make it work with I2P. This document concerns those changes and why they make sense for I2P Browsing,
Mixed I2P/Tor browsing, and I2P/Outproxy browsing.

The Basics:
-----------

The bare-minimum requirement to browse I2P using the Tor Browser is to add the following items to the
`user.js` file in your profile directory.

```javascript
user_pref("network.proxy.no_proxies_on", "127.0.0.1:7657,localhost:7657,127.0.0.1:7662,localhost:7662,127.0.0.1:7669,localhost:7669");
user_pref("extensions.torbutton.use_nontor_proxy", true);
user_pref("extensions.torlauncher.start_tor", false);
user_pref("extensions.torlauncher.prompt_at_startup", false);
user_pref("network.proxy.type", 1);
user_pref("network.proxy.http", "127.0.0.1");
user_pref("network.proxy.http_port", 4444);
user_pref("network.proxy.ssl", "127.0.0.1");
user_pref("network.proxy.ssl_port", 4444);
user_pref("network.proxy.ftp", "127.0.0.1");
user_pref("network.proxy.ftp_port", 4444);
user_pref("network.proxy.socks", "127.0.0.1");
user_pref("network.proxy.socks_port", 4444);
user_pref("network.proxy.share_proxy_settings", true);
user_pref("browser.startup.homepage", "about:blank");
```

However, to do so would then break your Tor Browser's ability to connect to Tor in normal, non-I2P cirumstances.
This is not desirable, so instead we place a `user.js` and `prefs.js` file in a fresh profile directory, in advance
of the user starting the Tor Browser. Which brings us to the most important thing:

### RULE ZERO: Preserve Normal Tor Browser Functionality when the user requests Tor Browser

If the user requests Tor Browser, we do not interfere in any way. We just start Tor Browser normally. I2P operation
is always separated to it's own profile.

Beside that, we make the remaining changes which are specified in the [i2p.firefox](https://i2pgit.org/i2p-hackers/i2p.firefox)
`user.js` and `prefs.js` file. This is a decision made to match the `i2p.firefox` profile as closely as possible,
however it may not be advisable. The question is "When running in I2P mode, are we trying to look like Tor Browser?"
If so, then we need to cut the user.js down to *just this section* but if not, then we should probably leave it alone.

IMO, it is not feasible to make our modifications to Tor Browser entirely invisible. I2P

Defending from Attacks Unique to Mixed Tor/I2P Environment:
-----------------------------------------------------------

The major risk with using I2P and Tor together in the same browser is that the identify of a user on one network
may leak across the networks, i.e. an onion site learns the identify of an I2P client tunnel or an I2P site learns
the identity of a Tor client. This can mean A) Cryptographic identity or B) ISP identity. ISP identity is largely
defended against by Tor Browser's extensive attention to proxy obedience, but we must be careful not to break this.
Cryptographic identity, as presented by `X-I2P-*` headers in I2P for example, is easier for an attacker to gain.

In order to defend against these kinds of attacks, when running in I2P mode this system adds several extensions to
the Tor Browser. In addition to NoScript and HTTPS Everywhere, which are offered by Tor Browser, in I2P mode the
following extensions will be added:

 1. I2P in Private Browsing
 2. Onion in Container Browsing
 3. JS Restrictor
 4. uBlock Origin
 5. LocalCDN - Fork of Decentraleyes
 and optionally,
 6. Actually Work Offline

#### Extensions? Isn't it bad to add extensions to Tor Browser?

> Well yes, if you're browsing Tor or the Web, and you're sure that you are starting with a common fingerprint.
However I2P has a *"fragmented"* fingerprint, and mixed Tor/I2P browsing creates an intrinsically *uncommon*
fingerprint. We can only benefit, in terms of security, from unifying upon a configuration(and thus a common
fingerprint).

> The danger in adding extensions to Tor Browser lies in the fact that you become a unique character among Tor
Browser users for using that extension. If *everyone* in your anonymity set, i.e. mixed Tor/I2P users uses the
same extensions when browsing I2P in Tor Browser, then you are not more unique for using the extensions.

> This experimental extension set is therefore designed to prevent known attacks on mixed I2P/Tor Browser users
and optimize Tor Browser running in Mixed Tor/I2P mode for I2P use. You should assume that it makes you look
like other users of this configuration, **not** like the Tor Browser.

### Protecting Network Boundaries:

I2P is normally administered via a WebUI(The "Router Console"), which the user sometimes views in the same
browser they use to visit remote sites. It is essential to prevent sites from being able to access information
from this WebUI and other locally running sites.

Clearnet sites will remain in the `firefox-default` container, but as soon as the firefox-default container
requests a `*.onion` or `*.i2p` site, it will be instantly containerized. Tor Browser's identity management
tooling controls clearnet tabs, and they are proxied using the default I2P outproxy or outproxy plugin.

#### Implemented in Extensions:

 - JS Restrictor: Prevents any non-local host from using Javascript to connect to the localhost.
 - I2P in Private Browsing: Places I2P browsing into it's own Container Tab where requests to `*.onion` resources
 are automatically dropped.
 - Onion in Container Tabs: Places .onion browsing into it's own Container Tab where requests to `*.i2p` resources
 are automatically dropped.

### Limiting Clearnet Accesses

I2P has been designed primarily with an in-network approach. It therefore makes sense to limit the use
of a "Backup" clearnet access like Tor or an outproxy. This should have the added benefit of an improved UX for
I2P browser users by blocking ads and speeding up access to resources that would normally be hosted by CDN's.

#### Implemented in Extensions

 - uBlock Origin: Block ads, and by extension much malvertising. Reduce the potential of an attacker using an
 advertising network to deliver an attack to I2P users. Prevent ad networks from using their reach to build
 profiles of I2P users. And there's no point serving ads to I2P users through an outproxy.
 - LocalCDN - Fork of Decentraleyes: LocalCDN is an extension in your browsing that acts as a cache of resources
 which would normally be provided by CDN's operated by companies that operate on the clearnet. Using LocalCDN
 for our browsing prevents a malfunctioning outproxy from affecting site functionality, decreases the load on
 outproxy operators, and improves UX by making sites load faster.
 - Actually Work Offline: It is possible to pass the `-offline` flag `-i2pconfig` flag to load an additonal
 extension wherein *all non-localhost requests* are immediately dropped and the browser cannot reach remote sites
 until restarted.

Other Reasons to put Extensions in Tor Browser for Our Purposes
---------------------------------------------------------------

 - Announce our non-default configuration loudly to the user by making changes to the UI.
 
* Unrelated to extensions: Additionally, apply the environment variable: `TOR_HIDE_BROWSER_LOGO=1` for all non-Tor
uses.