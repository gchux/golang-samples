/**
 *
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

const { HelloRequest, RepeatHelloRequest, HelloReply } = require('./helloworld_pb.js');
const { GreeterClient } = require('./helloworld_grpc_web_pb.js');

const params = new Proxy(new URLSearchParams(window.location.search), {
  get: (searchParams, prop) => searchParams.get(prop),
});
let server = params.server;

var client = new GreeterClient(`https://${server}`, null, null);

// simple unary call
var request = new HelloRequest();
request.setName('World');

client.sayHello(request, {}, (err, response) => {
  if (err) {
    console.log(`Unexpected error for sayHello: code = ${err.code}` +
      `, message = "${err.message}"`);
  } else {
    console.log(response.getMessage());
  }
});
