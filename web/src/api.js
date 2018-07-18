import axios from 'axios'

axios.interceptors.request.use(config => {
    const token = localStorage.getItem('Authorization')
    config.headers.common['Authorization'] = 'Bearer ' + token
    return config;
  },  error => {
    return Promise.reject(error);
  });

axios.interceptors.response.use(function (response) {
    if (response.data.error && response.data.error.code === 401) {
        location.replace('https://mixin.one/oauth/authorize?client_id=' + process.env.VUE_APP_CLIENT_ID + '&scope=PROFILE:READ&response_type=code')
        return Promise.reject(response)
    }
    return response;
  }, function (error) {
      if (error.response) {
          switch (error.response.status) {
            case 401:
          }
      }
    return Promise.reject(error);
  });

export default {
    getMe() {
        return axios.get("/me")
    },
    authenticate(authorizationCode) {
        var params = {
            "code": authorizationCode
        };
        return axios.post("/auth", params) 
    },
    getBot() {
        return axios.get("/bot")
    },
    verifyPayment(traceId) {
        return axios.get("/verify/"+traceId)
    },
    pushKeys(botId, clientId, sessionId, privateKey) {
        var params = {
            "client_id": clientId,
            "session_id": sessionId,
            "private_key": privateKey
        };
        return axios.post("/bot/"+botId+"/keys", params)
    }
}