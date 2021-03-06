const apiUrl = "https://yetian.design/ghost/api/v0.1"

//get Ghost token in browser's localStorage
var ls = localStorage.getItem("ghost:session")
//no token, redirect to Ghost Admin page
if (!ls) {
    localStorage.clear()
    window.location.href = "/ghost"
}
var ghostToken = JSON.parse(ls).authenticated.access_token

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
            slug: "",
            posts: [],
            barValue: 0,
            isShowBar: false
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
                                    <input class="file-input" type="file" accept="image/png" name="cover" @change="selectCover($event,p.slug)">
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
            <progress class="progress is-small" :value="barValue" max="100" v-show="isShowBar"></progress>
        </div>
    `,
    methods: {
        fetchData() {
            //get all posts
            axios
                .get(apiUrl + "/posts/", {
                    params: {
                        limit: 100,
                        fields: "title,slug,visibility"
                    }
                })
                .then(response => {
                    this.posts = response.data.posts
                })
        },
        selectCover: function(e, s) {
            this.cover = e.target.files[0]
            this.slug = s
            //check uploaded cover type
            if (
                this.cover.type != "image/png" ||
                this.cover.name.substr(-4) != ".png"
            ) {
                this.cover = null
                this.slug = ""
                alert("图片仅限PNG类型，请重新选择")
                return
            }
            this.uploadCover()
        },
        uploadCover: function() {
            this.isShowBar = true
            var sendData = new FormData()
            sendData.append("coverFile", this.cover)
            sendData.append("coverSlug", this.slug)
            var config = {
                onUploadProgress: progressEvent => {
                    this.barValue = Math.round(
                        (progressEvent.loaded * 100) / progressEvent.total
                    )
                }
            }
            axios
                .post("/upload/api/cover", sendData, config)
                .then(response => {
                    this.isShowBar = false
                    this.barValue = 0
                    if (response.data.status == "OK") alert("封面更换成功")
                    this.fetchData()
                })
                .catch(function(err) {
                    alert("上传失败：" + err.message)
                    console.error(err)
                    this.fetchData()
                })
        }
    },
    beforeMount: function() {
        this.fetchData()
    },
    beforeCreate: function() {
        //check token valid
        axios
            .get(apiUrl + "/users/1/")
            .then(response => {
                if (response.status != 200) {
                    localStorage.clear()
                    window.location.href = "/ghost"
                }
            })
            .catch(error => {
                localStorage.clear()
                window.location.href = "/ghost"
            })
    }
})

new Vue({
    el: "#cover-app"
})

//listen attachment uploader to handle upload event
var attInput = document.querySelector("#attachment-app input")
attInput.onchange = function(event) {
    var attachmentFile = event.target.files[0]
    var sendData = new FormData()
    sendData.append("attachment", attachmentFile)
    //upload progress
    var config = {
        onUploadProgress: function(progressEvent) {
            var percentCompleted = Math.round(
                (progressEvent.loaded * 100) / progressEvent.total
            )
            var bar = document.getElementById("att-progress-bar")
            bar.style.display = "block"
            bar.setAttribute("value", percentCompleted)
        }
    }
    axios
        .post("/upload/api/upload-attachment", sendData, config)
        .then(function(response) {
            var bar = document.getElementById("att-progress-bar")
            bar.style.display = "none"
            if (response.data.status == "OK") alert("上传成功")
            window.location.reload()
        })
        .catch(function(err) {
            alert("上传失败：" + err.message)
            console.error(err)
        })
}

//delete attachment handler
function deleteFile(target) {
    var filename = target.dataset.filename
    if (confirm("确认删除文件：" + filename)) {
        axios
            .post("/upload/api/delete-attachment", {
                filename: filename
            })
            .then(function(response) {
                if (response.data.status == "OK") alert("删除成功")
                window.location.reload()
            })
            .catch(function(err) {
                alert("删除失败：" + err.message)
                console.error(err)
            })
    } else {
        return
    }
}

//copy link button handler
var cb = new ClipboardJS(".copy-link-btn")
cb.on("success", function(e) {
    alert("已复制文件地址：" + e.text)
})
cb.on("error", function(e) {
    alert("浏览器不支持复制操作")
})
