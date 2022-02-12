package tbsupervise

var secbrowserjs = []byte(`
// Copyright (C) 2019 - 2021 ENCRYPTED SUPPORT LP <adrelanos@whonix.org>
// See the file COPYING for copying conditions.

// Warning! These settings disable Tor. You will not be anonymous!

// Configure Tor Browser without Tor settings for an everyday use
// security hardened browser. Take advantage of its excellent
// enhancements for reducing linkability, that is, "the ability
// for a user's activity on one site to be linked with their
// activity on another site without their knowledge or explicit
// consent."
// - See https://www.whonix.org/wiki/Tor_Browser_without_Tor
// - See https://www.whonix.org/wiki/SecBrowser

// This file gets copied at first start to:
// ~/.secbrowser/secbrowser/Browser/TorBrowser/Data/Browser/profile.default/user.js

// Disable Torbutton and Torlauncher extensions
user_pref("extensions.torbutton.startup", false);
user_pref("extensions.torlauncher.start_tor", false);
user_pref("network.proxy.socks_remote_dns", false);

// https://forums.whonix.org/t/tor-browser-10-without-tor/10313
user_pref("network.dns.disabled", false);

// Disable function torbutton source file:
// src/chrome/content/torbutton.js
// function: torbutton_do_tor_check
// and thereby also function: torbutton_initiate_remote_tor_check
// disables Control Port verification and remote Tor connection check.
user_pref("extensions.torbutton.test_enabled", false);

// Set security slider "Safest".
// Broken, therefore disabled by default.
// https://forums.whonix.org/t/broken-security-slider/8675
// user_pref("extensions.torbutton.inserted_security_level", true);
// user_pref("extensions.torbutton.security_slider", 1);

// Normalize Tor Browser behavior
user_pref("extensions.torbutton.noscript_persist", true);
user_pref("browser.privatebrowsing.autostart", false);

// Save passwords.
user_pref("signon.rememberSignons", true);

// Required for saving passwords.
// https://trac.torproject.org/projects/tor/ticket/30565#comment:7
user_pref("security.nocertdb", false);

// Disable Letterboxing.
// https://forums.whonix.org/t/is-anyone-having-white-bars-in-the-tbb-tor-browser-letterboxing/8345
// https://forums.whonix.org/t/secbrowser-a-security-hardened-non-anonymous-browser/3822/156
user_pref("privacy.resistFingerprinting.letterboxing", false);

// Enable punycode to fix
// very hard to notice Phishing Scam - Firefox / Tor Browser URL not showing real Domain Name - Homograph attack (Punycode).
// https://forums.whonix.org/t/very-hard-to-notice-phishing-scam-firefox-tor-browser-url-not-showing-real-domain-name-homograph-attack-punycode/8373
// https://forums.whonix.org/t/secbrowser-a-security-hardened-non-anonymous-browser/3822/162
user_pref("network.IDN_show_punycode", true);

// Disable popup asking to prefer onions since onions will not work in clearnet browser.
user_pref("privacy.prioritizeonions.showNotification", false);
`)

var secbrowserhtml = []byte(`
<!DOCTYPE html>
<html>
<head>
  <title>i2p.plugins.tor-manager - Tor Binary Manager</title>
  <link rel="stylesheet" type="text/css" href ="/style.css" />
</head>
<body>
<h1 id="running-in-clearnet-mode">Running in Clearnet Mode</h1>
<p>Tor Browser is configured to run without Tor, and will now use the non-anonymous web. It is also configured to use uBlock Origin. This allows you to use a hardened web browser for your non-anonymous tasks.</p>
<ul>
<li>To get started, perhaps try <a href="https://duckduckgo.com">DuckDuckGo</a></li>
<li>or <a href="https://privacyguides.org/">PrivacyGuides</a></li>
</ul>
<p>This wrapper has been developed for use with the I2P project. To learn more about I2P, visit https://geti2p.net</p>
</body>
</html>
`)

var offlinehtml = []byte(`
<!DOCTYPE html>
<html>
<head>
  <title>i2p.plugins.tor-manager - Tor Binary Manager</title>
  <link rel="stylesheet" type="text/css" href ="/style.css" />
</head>
<body>
<h1 id="running-in-offline-mode">Running in Offline Mode</h1>
<p>Tor Browser is configured to run without <em>any</em> access to the clearnet, and will now only use local services running on this device. This uses a different mechanism than Firefox’s normal work offline mode and cannot be canceled.</p>
<ul>
<li>Visit the <a href="http://127.0.0.1:7657">I2P Router Console</a></li>
<li>Visit the <a href="http://127.0.0.1:8888">Freenet FProxy</a></li>
</ul>
<p>This wrapper has been developed for use with the I2P project. To learn more about I2P, visit https://geti2p.net</p>
</body>
</html>
`)

var defaultCSS []byte = []byte(`
* {
	padding: 0;
	margin: 0;
  }
  
  html {
	margin: 0 4%;
	padding: 0 20px;
	min-height: 100%;
	background: #9ab;
	background: repeating-linear-gradient(to bottom, #9ab, #89a 2px);
	scrollbar-color: #bcd #789;
  }
  
  body {
	margin: 0;
	padding: 20px 40px;
	font-family: Open Sans, Noto Sans, Segoe UI, sans-serif;
	font-size: 12pt;
	color: #495057;
	text-decoration: none;
	word-wrap: break-word;
	border-left: 1px solid #495057;
	border-right: 1px solid #495057;
	box-shadow: 0 0 2px 2px rgba(0, 0, 0, .1);
	background: #f2f2f2;
  }
  
  h1, h2, h3, h4 {
	display: block;
	font-weight: 700;
  }
  
  h1 {
	text-transform: uppercase;
	font-weight: 900;
	font-size: 200%;
  }
  
  h2 {
	font-size: 140%;
  }
  
  h3 {
	font-size: 120%;
  }
  
  h4 {
	margin-bottom: 5px;
	text-align: right;
	text-transform: none;
	font-size: 90%;
	font-weight: 600;
	font-style: italic;
  }
  
  p {
	margin-bottom: 15px;
	width: 100%;
	line-height: 1.4;
	word-wrap: break-word;
  /*  text-align: justify;*/
	text-decoration: none;
  }
  
  ul {
	margin: 10px 20px;
	list-style: none;
  }
  
  li {
	margin-left: 0;
	padding: 12px 15px 15px 20px;
	width: calc(100% - 40px);
	text-align: justify;
	border: 1px solid #9ab;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff;
	background: #dee2e6;
  }
  
  li li {
	padding-bottom: 0;
	width: calc(100% - 40px);
	text-align: left;
	border: none;
	border-top: 1px solid #9ab;
	box-shadow: none;
  }
  
  li li:first-of-type {
	margin-top: 15px;
	border-top: none;
  }
  
  li a:first-of-type {
	display: block;
	width: 100%;
  }
  
  #applicationExplain {
	float: unset;
  }
  
  li+li {
	margin-top: 15px;
  }
  
  h3+ul, ul+h3, ul+h2 {
	margin-top: 20px;
  }
  
  a, button {
	color: #3b6bbf;
	text-decoration: none;
	font-weight: 700;
	word-wrap: break-word;
	outline: 0;
  }
  
  .applicationDesc {
	color: #81888f;
	text-decoration: none;
	font-weight: 700;
	word-wrap: break-word;
	outline: 0;
  }
  
  .applicationDesc:hover, a:hover, button:hover {
	text-decoration: none;
	font-weight: 700;
	word-wrap: break-word;
	outline: 0;
  }
  
  button {
	border: none;
	cursor: pointer;
	color: #3b6bbf;
	text-decoration: none;
	font-weight: 700;
	word-wrap: break-word;
	outline: 0;
  }
  
  .background {
	background-color: #f8f8ff;
	height: 100%;
  }
  
  .content {
	margin: 1.5rem;
	padding: 1rem;
	min-height: 3rem;
	min-width: 95%;
	display: inline-block;
	border: 1px solid #d9d9d6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
	background: #f8f8ff;
  }
  
  #header, .application-info, .browser-info, .extended-info, .search-info {
	margin-top: 1.5rem;
	padding: 1rem;
	min-height: 3rem;
	min-width: 95%;
	display: inline-block;
	border: 1px solid #d9d9d6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
	background: #f8f8ff;
  }
  
  .showhider {
	margin-right: auto;
	padding: 0!important;
	text-transform: uppercase;
	background: none !important;
	border: none;
	width: 90%;
	color: #3b6bbf;
	text-decoration: none;
	font-weight: 700;
	word-wrap: break-word;
	outline: 0;
	text-align: left;
  }
  
  #links .showhider {
	font-size: 25px;
  }
  
  .section-header {
	display: flex;
	flex-direction: row;
	margin-bottom: 80px;
  }
  
  #readyness {
	padding-top: 1rem;
	padding-bottom: 1rem;
	margin: 1rem;
	width: 42%;
	min-width: 42%;
	background: #dee2e6;
	text-align: center!important;
	border: 1px solid #dee2e6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
	display: inline-block;
  }
  
  #onboarding {
	min-height: 5rem;
	padding: .5rem;
	margin: .5rem;
	margin-top: 4rem;
	width: 42%;
	min-width: 42%;
	font-size: 2rem;
	background: #a48fe1;
	text-align: center!important;
	border: 1px solid #a48fe1;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  #i2pbrowser-description {
	padding-top: 1rem;
	padding-bottom: 1rem;
	width: 50%;
	min-width: 50%;
	display: inline-block;
	background: #dee2e6;
	border: 1px solid #dee2e6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  #linksExplain {
	min-height: 5rem;
	padding: .5rem;
	margin: .5rem;
	width: 30%;
	min-width: 30%;
	background: #dee2e6;
	text-align: center!important;
	border: 1px solid #dee2e6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  #applicationExplain, #controlExplain {
	min-height: 5rem;
	padding: .5rem;
	margin: .5rem;
	width: 30%;
	min-width: 30%;
	background: #dee2e6;
	text-align: center!important;
	border: 1px solid #dee2e6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
	float: left;
  }
  
  #proxyReady {
	min-height: 3rem;
	padding: .5rem;
	margin: .2rem;
	width: 38%;
	min-width: 38%;
	display: inline-block;
	background: #d9d9d6;
	float: right;
	text-align: center!important;
	border: 1px solid #d9d9d6;
	border-radius: 2px;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  #proxyUnready {
	min-height: 3rem;
	padding: .5rem;
	margin: .2rem;
	width: 38%;
	min-width: 38%;
	display: inline-block;
	float: right;
	text-align: center!important;
	border: 1px solid #ffc56d;
	border-radius: 2px;
	background: #ffc56d;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  #consoleOn {
	min-height: 3rem;
	padding: .5rem;
	margin: .2rem;
	width: 38%;
	min-width: 38%;
	display: inline-block;
	float: left;
	text-align: center!important;
	border: 1px solid #f7e59a;
	border-radius: 2px;
	background: #f7e59a;
	box-shadow: inset 0 0 0 1px #fff, 0 0 1px #ccc;
  }
  
  .onboardingContent {
	font-size: .8rem!important;
	text-align: left;
	display: none;
  }
  
  #info-content {
	display: none;
  }
  
  .consoleOn:hover #proxy-check, .proxyReady:hover #proxy-check {
	visibility: visible;
	opacity: 1;
  }
  
  img {
	max-width: 100%;
  }
  
  img.readyness {
	height: 100%;
	width: auto;
  }
  
  @media only screen and (max-width: 399px) {
	.application-info {
	  display: none;
	}
  }
  
  @media screen and (max-width: 1200px) {
	body {
	  font-size: 10.5pt;
	}
  }
  
  video {
	width: 100%
  }
  
   /* The switch - the box around the slider */
.switch {
	position: relative;
	display: inline-block;
	width: 60px;
	height: 34px;
  }
  
  /* Hide default HTML checkbox */
  .switch input {
	opacity: 0;
	width: 0;
	height: 0;
  }
  
  /* The slider */
  .slider {
	position: absolute;
	cursor: pointer;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: #ccc;
	-webkit-transition: .4s;
	transition: .4s;
  }
  
  .slider:before {
	position: absolute;
	content: "";
	height: 26px;
	width: 26px;
	left: 4px;
	bottom: 4px;
	background-color: white;
	-webkit-transition: .4s;
	transition: .4s;
  }
  
  input:checked + .slider {
	background-color: #2196F3;
  }
  
  input:focus + .slider {
	box-shadow: 0 0 1px #2196F3;
  }
  
  input:checked + .slider:before {
	-webkit-transform: translateX(26px);
	-ms-transform: translateX(26px);
	transform: translateX(26px);
  }
  
  /* Rounded sliders */
  .slider.round {
	border-radius: 34px;
  }
  
  .slider.round:before {
	border-radius: 50%;
  } 

`)
