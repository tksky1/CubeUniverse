type methods = "get" | "post";

type apis = {
    "auth/login": {
        method: "get",
        params: {
            uid: string,
            password: string
        }
    },
    "auth/altpas": {
        method: "post"
        params: {
            uid: string,
            password: string,
            newpassword: string,
            newuid: string
        }
    }
}

export default function api<T extends keyof apis>(route: T, params: apis[T]) {
}
