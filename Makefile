gen:
    # Generate twirp code
	protoc --go_out=. --twirp_out=. rpc/charlie-go/service.proto

