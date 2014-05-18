
(function() {
    var user = localStorage["user"];
    if(!user) {
        promptAndLogin();
    } else {
        login();
    }


    function promptAndLogin() {
        user = prompt("who are you?");
        if(!user) {
            promptAndLogin();
        } else {
            localStorage["user"] = user = user.trim();
            login();
        }
    }

    function login() {
        loadAndShowStatuses();
        sup();
    }

    function loadAndShowStatuses(){
        d3.json("api/status/"+user)
            .get(function(error, data){
                if(!error){
                    d3.select(".js-sup").text("sup "+user+"?");
                    showStatuses(data);
                } else if(confirm("couldn't load dc :( try again?")){
                    login();
                }
            });
    }

    function showStatuses(data) {
        console.log("Displaying "+data.length+" statuses");
        var statusHtml = function(d){
            return "<span class='text-muted'>"+d.CreateTime+"</span> " + 
                d.Tags.split(" ").map(function(tag){
                    return "<span class='label label-default'>"+tag+"</span>";
                }).join(" ") + " " +
                d.Description;
        };
        var d = d3.select(".js-statuses").selectAll("div")
            .data(data)
            .html(statusHtml);
        d.enter().append("div")
            .html(statusHtml);
        d.exit().remove();
    }

    function sup() {
        var input = prompt("what are you doing right now?\nformat: 'tag1 tag2 ...: desc'");
        if(input != null){
            post(parse(input));
        }

        // Wait a random amount of time. 
        // At least 1 hour, at most 150 minutes = 2.5 hours. Usually slightly over 1h.
        var waitMins = 60;
        for(var prob = 0.1; Math.random() > prob; prob += 0.1) {
            waitMins += 10;
        }
        console.log("Sleeping for "+waitMins+" minutes");
        window.setTimeout(sup, waitMins*60*1000);
    }

    function parse(input){
        var tags, desc;
        var ix = input.indexOf(":");
        if(ix > -1) {
            tags = parseTags(input.substring(0,ix));
            desc = input.substring(ix+1).trim();
        } else {
            tags = parseTags(input);
            desc = "";
        }
        return {
            "user": user,
            "tags": tags.join(" "),
            "description": desc
        };
    }

    function parseTags(input){
        return input.trim().split(/\s+/);
    }

    function post(body){
        d3.xhr("api/status")
            .header("Content-Type", "application/json")
            .post(JSON.stringify(body), function(error, data){
                if(!error) {
                    console.log("Posted successfully!");
                    loadAndShowStatuses();
                } else if(confirm("oh no, post failed. try again?")){ 
                    post(body);
                }
            });
    }
})();
