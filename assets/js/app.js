const apiUrl = "https://yetian.design/ghost/api/v0.1"

//get Ghost token in browser's localStorage
var ls = JSON.parse(localStorage.getItem("ghost:session"))
//no token, redirect to Ghost Admin page
if (!ls) {
    window.location.href = "/ghost"
}
var ghostToken = ls.authenticated.access_token

//inject axios into Vue.prototype
window.axios = axios
//inject Ghost token into axios header
if (ghostToken) {
    window.axios.defaults.headers.common["Authorization"] =
        "Bearer " + ghostToken
}

Vue.component("ghost-covers", {
    data: function() {
        return {
            message: "Covers",
            cover: null,
            posts: []
        }
    },
    template: `
        <div>
            <h1 class="title">{{ message }}</h1>
            <table class="table is-bordered">
                <thead>
                    <th>Title</th>
                    <th>Cover</th>
                    <th>Visibility</th>
                    <th>Operation</th>
                </thead>
                <tbody>
                    <tr v-for="p in posts">
                        <td>{{p.title}}</td>
                        <td><img :src="'/content/images/cover/'+p.slug+'.png'" :alt="p.slug" style="max-width:100px"></td>
                        <td>{{p.visibility}}</td>
                        <td>
                            <div class="file" style="margin-bottom:15px;">
                                <label class="file-label">
                                    <input class="file-input" type="file" accept="image/png" name="cover" @change="selectCover">
                                    <span class="file-cta button is-small is-outlined">
                                        <span class="file-label">
                                            Change
                                        </span>
                                    </span>
                                </label>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    `,
    methods: {
        selectCover: function(e) {
            this.cover = e.target.files[0]
        }
    },
    beforeMount: function() {
        //get all posts
        var self = this
        axios
            .get(apiUrl + "/posts", {
                params: {
                    limit: 100,
                    fields: "title,slug,visibility"
                }
            })
            .then(function(response) {
                self.posts = response.data.posts
            })
    }
})

new Vue({
    el: "#cover-app"
})

document.addEventListener("DOMContentLoaded", function() {
    //listen attachment uploader to handle upload event
    var attInput = document.querySelector("#attachment-app input")
    attInput.onchange = function(event) {
        console.log(event.target.files[0])
    }
})
