If you have an occasional need to share files between a variety of devices (ie. Windows machine, Linux laptop, iOS tablet), and you don't want to set up a heavier solution like OwnCloud or similar, this project may be of use to you.

* No Go dependencies
* Upload and download files using only a web browser
* Easy to build, run
* Striving for low/no maintenance

```
go build main.go
./main
```

The app will listen on the port defined in the code. Visit the URL in a browser and there you can upload files. Once a file is uploaded, a link to download it shows on the homepage. That's it.
