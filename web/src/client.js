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

const id = params.id || crypto.randomUUID();
const host = params.host;
const port = params.port;
const name = params.name;
const token = params.token;

class GrpcClientProvider {

  // client takes 3 arguments to instantiate
  constructor(client, host, port, id) {
    this.client = client;
    this.defaultServer = `${host}:${port}`;
    this.clientID = id;
  }

  mapMetadata(client, token) {
    for (const prop in client) {
      if (typeof client[prop] !== 'function') {
        continue
      }

      const original = client[prop]
      client[prop] = ((...args) => {
        args[1] = {
          ...args[1],
          Authorization: `Bearer ${token}`,
          'x-client-id': this.clientID,
        };
        return original.call(client, ...args);
      })
    }
    return client;
  }

  clientForServer(server) {
    return new this.client(server, null, null);
  }

  defaultClient() {
    return this.clientForServer(this.defaultServer, null, null);
  }

  authorizedClient(token) {
    return this.mapMetadata(this.defaultClient(), token);
  }

  authorizedClientForServer(server, token) {
    return this.mapMetadata(this.clientForServer(server), token)
  }

  id() {
    return this.clientID;
  }
}

const exec = function (client, request) {
  client.sayHello(request, {}, (err, response) => {
    if (err) {
      console.log(id, `Unexpected error for sayHello: code = ${err.code}` +
        `, message = "${err.message}"`);
    } else {
      console.log(id, response.getMessage());
    }
  });

}

const newRequest = function (name) {
  // simple unary call
  const request = new HelloRequest();
  request.setName(name);
  return request;
}

const clientProvider = new GrpcClientProvider(GreeterClient, host, port, id)
const authorizedClient = token ? clientProvider.authorizedClient(token) : clientProvider.defaultClient()

exec(authorizedClient, newRequest(name))

window['grpc_client'] = {
  get: () => clientProvider.defaultClient(),
  sayHello: (name) => exec(clientProvider.defaultClient(), newRequest(name)),
  sayHello_withToken: (token) => exec(clientProvider.authorizedClient(token), newRequest(name)),
  sayHello_withServerAndToken: (host, port, token) => exec(clientProvider.authorizedClientForServer(`${host}:${port}`, token), newRequest(name)),
}
