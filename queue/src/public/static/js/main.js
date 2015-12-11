/**
 * Created by yaoyun on 15/12/10.
 */
function Flush() {
    $.ajax({
        type: "GET",
        url: "/v1/events/1",
        data: {},
        success: function(data){
            console.log(data["data"])
            if (data["data"]["unlockOn"] > data["data"]["curTime"]){
                timer.resetCountdown((data["data"]["unlockOn"]- data["data"]["curTime"]) * 100)
                $("#countdown").show()
                $("#btWait").show()
                $("#btOver").hide()
                $("#btBuy").hide()
            }
            else {
                $("#countdown").hide()
                $("#btWait").hide()
                $("#btOver").show()
                $("#btBuy").hide()
            }
        }
    });
}


function Buy() {
    $.ajax({
        type: "POST",
        url: "/v1/tickets",
        crossDomain: true,
        dataType: "json",
        data: {},
        success: function(data) {
            console.log(data);
            if (data["code"] == 0) {
                location.href = "/view/index-success.html"
                alert("Congratulation !!! You Succeed !!!")
            }
            else {
                alert("Game Over")
                $("#countdown").hide()
                $("#btWait").hide()
                $("#btOver").show()
                $("#btBuy").hide()
            }
        }
    });
}

$(document).ready(function(){
    Flush()
});

// Common functions
function pad(number, length) {
    var str = '' + number;
    while (str.length < length) {str = '0' + str;}
    return str;
}

function formatTime(time) {
    var min = parseInt(time / 6000),
        sec = parseInt(time / 100) - (min * 60),
        hundredths = pad(time - (sec * 100) - (min * 6000), 2);
    return (min > 0 ? pad(min, 2) : "00") + ":" + pad(sec, 2) + ":" + hundredths;
}

var timer = new (function() {
    var $countdown,
        incrementTime = 70,
        currentTime = 300000000,
        updateTimer = function() {
            $countdown.html("倒计时: " + formatTime(currentTime));
            if (currentTime == 0) {
                timer.Timer.stop();
                timerComplete();
                return;
            }
            currentTime -= incrementTime / 10;
            if (currentTime < 0) currentTime = 0;
        },
        timerComplete = function() {
            $("#btWait").hide()
            $("#btOver").hide()
            $("#btBuy").show()
            alert('活动开始');
        },
        init = function() {
            $countdown = $('#countdown');
            timer.Timer = $.timer(updateTimer, incrementTime, true);
        };
    this.resetCountdown = function(newTime) {
        if (newTime > 0) {currentTime = newTime;}
    };
    $(init);
});

$("#btBuy").click(function() {
    Buy()
});
