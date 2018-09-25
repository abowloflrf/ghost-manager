const apiUrl = "https://yetian.design/ghost/api/v0.1"

//get Ghost token in browser's localStorage
var ls = JSON.parse(localStorage.getItem("ghost:session"))
var ghostToken = ls.authenticated.access_token

//inject axios into Vue.prototype
Vue.prototype.$axios = axios
//inject Ghost token into axios header
if (ghostToken) {
    Vue.prototype.$axios.defaults.headers.common["Authorization"] =
        "Bearer" + ghostToken
} else {
    window.location.href = "/ghost"
}

Vue.component("ghost-covers", {
    data: function() {
        return {
            message: "Hello Vue"
        }
    },
    template: "<p>{{ message }}</p>",

    beforeMount: function() {
        this.$axios.get(apiUrl + "/post", {
            params: {
                limit: 100,
                fields: "title,slug,visibility"
            }
        }).then(function(response){
            console.log(response.data)
        })
    }
})

new Vue({
    el: "#cover-app"
})
