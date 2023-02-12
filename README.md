# meta_dumper

This small script dynamically dumps all nodedefs and itemdefs into json format

### How to use?

1. download `git clone https://github.com/ev2-1/meta_dumper`
2. compile `cd meta_dumper; go build`
3. start the minetest server you want to dump meta from
   write down the port
4. start the dumper `./meta_dumper minetest:port ":freeport"`
5. connect to ":freeport" using a normal minetest client
6. DO NOT PANIC if it crashes, thats expected
7. done, you should now have a nodemeta.json and itemmeta.json file
