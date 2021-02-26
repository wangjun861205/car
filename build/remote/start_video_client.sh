#!/bin/sh

nc 192.168.3.8 5000 | mplayer -noconsolecontrols -nosound -framedrop -x 640 -y 480 -fps 60 -demuxer +h264es -cache 2048 -
