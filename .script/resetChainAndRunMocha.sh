.script/shutdown.sh
sleep 4
.script/setup.sh
flag1="-b"
flag2="block"
.script/startup.sh $flag1 $flag2
cd .mocha
npm run test:awesome
