const http = require("http");
const fs = require("fs");
const path = require("path");

console.log("Hello World");
fs.readFileSync("video.mp4", (err, data) => {
  console.log("opened");
  if (err) {
    console.log(err);
  } else {
    console.log(data);
  }
});

const data = fs.readFileSync("./video.mp4");
console.log(data);

// console.log(Buffer.from(data).toString("base64"));

const options = {
  hostname: "localhost",
  port: 8080,
  path: "/utubeapi/v1/video/chunk",
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
};

console.log("Sending");
const req = http.request(options, (res) => {
  console.log(`statusCode: ${res.statusCode}`);

  res.on("data", (d) => {
    console.log(d.toString());
    // process.stdout.write(d)
  });
});

req.write(
  JSON.stringify({
    name: "insomnia_video_create",
    data: Buffer.from(data).toString("base64"),
  }),
);
req.end();
