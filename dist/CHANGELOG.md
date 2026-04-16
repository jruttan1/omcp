## Changelog
* f1ee712c6e5740f17dfcb479beffc78f1d31318a Add GitHub Actions workflow for Go project
* 590dcd808aa58d500fd2adcb3b3c0f398a9e6eb3 Fix HTML entity for arrow in README
* a2b7d40dac4ef4726769cc5ab18b829d1adeb062 Merge branch 'main' of https://github.com/jruttan1/omcp
* 38826c27e786e4a57e4353acf8eb7d6ec4d125fc chore: cleanup and polish
* 698666d79997b8f323c42518ca7c72c4a06b3a3e chore: delete test file
* 02059d9e42357031a3940647e09b181cbf387423 chore: remove compiled binary from tracking
* 53fd6a7712a8a5e28042661f1d80fa46ae54c863 feat: add basic server with some metadata
* 03b585f5c71a3eecbfad9d87ffe75ec304e8ea47 feat: add request body for post and put
* 12e77277bb144a6f7b21768c181098d5ecbae77b feat: add tool builder function, requestURL helper and response handler for server to agent communication
* ea3f4c63474e4ab71ca9e35c29d363b1f76190b6 feat: add tui and styling, integrate with main
* cb61ff134adf2c524afb046185f914f2227ff432 feat: build core logic for parsing OpenAPI yaml file into MCP tools
* b438c46318492222ce966b095fb69d291cc04268 feat: build full tui for file/url input and endpoint selection
* b7a2f07af9b53158b13d6768c8e7049f3c385df9 feat: parse params from yaml and add to tools
* 3bbc25214c7ea11871038c6526431fdd19730529 feat: split commands into init for setup and config, then secondary command for running server (to be called by agent)
* 7c973429e7c7f78b636d9afc20dd7cfc021dfc3b feat: ui polish and handling for api keys
* d72093ec78369f92c7f5e930eb2501f2182c201f first commit
* 53beee549e9a220f48e2cf1732fc32c54fcba5e8 fix: lower case method name causes request rejections
* c8a2c79c4a453923c85a7c1a839bdb630a731c29 fix: ui bugs with filtering state
* ee0d5166085c1fffd8a53f7865a0262063029a1c fix: update go version in github workflow
* fd873e44175e6f6c2f4ee9114b09d78a27d888fe refactor: clean switch logic for tui state to make it easier to understand
* d0c83aab72f06013f5adec6f0bfd5076352b7e5a refactor: split main into parse function and add template for TUI and mcp generator
