It's rare that the server with which we are communicating with always behaves as expected. 
In reality, it's not just the server, but any of the other networking devices that our application's request passes through may not behave optimally.
So, how does our client fare then?


Keeping that in mind, in here we enforce:
- Time-outs in our clients
- Create client middleware.
- Explore connection pooling.