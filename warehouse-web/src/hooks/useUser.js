export default function useUser() {
    const getUserId = () => {
        const USERID = localStorage.getItem('USERID')
        if (!USERID) {
            localStorage.removeItem('USERID')
            localStorage.removeItem('USERNAME')
            localStorage.removeItem('token')
        }
        return USERID
    }
    const setUserInfo = (response) => {
        localStorage.setItem('USERID', response.data.USERID)
        localStorage.setItem('USERNAME', response.data.USERNAME)
        localStorage.setItem('token', response.data.token)
    }
    const isAuthenticated = () => {
        const hasToken = localStorage.getItem('token')
        return hasToken
    }
    return { getUserId, setUserInfo, isAuthenticated }
}
