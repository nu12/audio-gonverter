angular.module('gonverter', [])
  .controller('GonverterController', function() {
    var c = this;
    c.message = "Welcome to augio-gonverter!"
    c.files = [
        {
            original: "file1.mp3",
            converted: "file1.ogg",
            status: "Ready",
            download: "abc"
        },
        {
            original: "file2.mp3",
            converted: "file2.ogg",
            status: "Converting",
            download: "#"
        },
        {
            original: "file3.mp3",
            converted: "file3.ogg",
            status: "Pending",
            download: "#"
        }
    ];
    c.upload = function() {
        c.message = "Uploading files..."
        filesToUpload = document.getElementById("files").files;
        console.log(filesToUpload);
        c.files = []
    };
  });