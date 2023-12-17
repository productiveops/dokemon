import axios from "axios"

axios.interceptors.response.use(
  (res) => {
    return res
  },
  (err) => {
    if (axios.isAxiosError(err)) {
      if (err.response?.status === 401) {
        window.location.pathname = "/login"
      }
    }
  }
)
