# clear /media/uploads/1
rm -rf /media/uploads/*

# copy /home/my-user/core/fixtures/media/video-test/* to /media/uploads/1/*
mkdir -p /media/uploads/1 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/1
mkdir -p /media/uploads/2 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/2
mkdir -p /media/uploads/3 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/3
mkdir -p /media/uploads/4 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/4
mkdir -p /media/uploads/5 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/5
mkdir -p /media/uploads/6 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/6
mkdir -p /media/uploads/7 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/7
mkdir -p /media/uploads/8 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/8
mkdir -p /media/uploads/9 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/9
mkdir -p /media/uploads/10 && cp -rf /usr/src/app/core/fixtures/media/video-test/* /media/uploads/10

# copy /home/my-user/core/fixtures/media/thumbnails/* to /media/uploads/* without create a folder
cp -r /usr/src/app/core/fixtures/media/thumbnails/* /media/uploads/
