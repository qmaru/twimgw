export namespace core {
	
	export class TwitterDownload {
	    username: string;
	    max_results: string;
	    start_id: string;
	    exclude: boolean;
	    socks5: string;
	    storage_path: string;
	
	    static createFrom(source: any = {}) {
	        return new TwitterDownload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.max_results = source["max_results"];
	        this.start_id = source["start_id"];
	        this.exclude = source["exclude"];
	        this.socks5 = source["socks5"];
	        this.storage_path = source["storage_path"];
	    }
	}

}

