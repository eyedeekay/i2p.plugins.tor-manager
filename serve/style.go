package tbserve

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
  
`)
