const fs = require("fs");
const http = require("http");
const HttpDispatcher = require("httpdispatcher");

const dispatcher = new HttpDispatcher();
const port = 3000;
let count = 0;

const contents = fs.readFileSync("big.txt", "utf8");
const fileNames = contents.split("\n");

const offset =
  Math.round(Math.random() * fileNames.length) % (fileNames.length - 100);

const requestHandler = (request, response) => {
  count++;
  console.log(count);

  try {
    console.log(`ID: ${count}: ${request.url}`);
    dispatcher.dispatch(request, response);
  } catch (err) {
    console.log(err);
  }
};

const validFileNames = fileNames.slice(offset, offset + 100);
for (const fileName of validFileNames) {
  dispatcher.onGet(`/${fileName}`, function (_, res) {
    res.writeHead(200, { "Content-Type": "text/html" });
    res.end("<h1>Hey, this is the homepage of your server</h1>");
  });
}

dispatcher.onError(function (_, res) {
  res.writeHead(404);
  res.end("Error");
});

const server = http.createServer(requestHandler);
// server.timeout = 90000;

server.listen(port, (err) => {
  if (err) {
    return console.log("Something bad happened", err);
  }

  console.log(`Server is listening on ${port} 🚀`);
  console.log(`File offset: ${offset}`);
  console.log(`Valid directory paths: ${validFileNames.length}`);
});
