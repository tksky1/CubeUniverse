import { authentication } from "@/storage";

function addAuthentication(headers: Headers) {
    headers.append("Authorization", authentication.value);
    // headers.append(
    //     "Authorization",
    //     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTY4MjE1NzgwNiwiaWF0IjoxNjgxNTUzMDA2fQ.oKyrHMEe3DRAfpRIKPG_bCe-hycVeLHeU-PU1kTbBf0");
}

export function checkStorageOpen(type: string) {
    let headers = new Headers();
    addAuthentication(headers);
    headers.append("X-action", "check");
    headers.append("X-type", type);
    return fetch("joe/api/storage/check", {
        method: "GET",
        headers
    });
}

export function openStorage(type: string) {
    let headers = new Headers();
    addAuthentication(headers);
    headers.append("X-action", "create");
    headers.append("X-type", type);
    return fetch("joe/api/storage/check", {
        method: "GET",
        headers
    });
}

export function getPvcInfo(type: string) {
    let headers = new Headers();
    addAuthentication(headers);
    headers.append("X-type", type);
    return fetch("joe/api/storage/pvc", {
        method: "GET",
        headers
    });
}

interface updatePvcVolumeInterface {
    name: string,
    namespace: string,
    volume: string,
    "X-type": string,
}
export function updatePvcVolume(props: updatePvcVolumeInterface) {
    let headers = new Headers();
    addAuthentication(headers);
    let formData = new FormData();
    (Object.keys(props) as Array<keyof typeof props>).forEach(key => {
        formData.append(key, props[key]!);
    });

    return fetch("joe/api/storage/pvcpatch", {
        method: "POST",
        headers,
        body: formData
    });
}

interface deletePvcInterface {
    name: string,
    namespace: string,
    "X-type": string,
}
export function deletePvc(props: deletePvcInterface) {
    let { name, namespace, "X-type": xType } = props
    let headers = new Headers();
    addAuthentication(headers);
    let formData = new FormData();
    (Object.keys(props) as Array<keyof typeof props>).forEach(key => {
        headers.append(key, props[key]);
        // formData.append(key, props[key]!);
    });

    // return fetch(`joe/storage/pvc?name=${name}&namespace=${namespace}&X-type=${xType}`, {
    return fetch(`joe/api/storage/pvc`, {
        method: "DELETE",
        headers,
        body: JSON.stringify(props),
    });
}

interface createPvcInterface {
    name: string,
    namespace: string,
    volume?: string,
    "X-type": string,
    autoscale: "true" | "false",
    maxobject?: string,
    maxgbsize?: string
}

export function createPvc(props: createPvcInterface) {
    let headers = new Headers();
    addAuthentication(headers);
    let formData = new FormData();
    (Object.keys(props) as Array<keyof typeof props>).forEach(key => {
        formData.append(key, props[key]!);
    });

    return fetch("joe/api/storage/pvc", {
        method: "POST",
        headers,
        body: formData
    });
}

export function login(username: string, password: string) {
    return fetch("joe/api/auth/login", {
        method: "POST",
        body: JSON.stringify({
            "uid": username,
            "password": password
            // "uid": "12345678901",
            // "password": "cubeuniverse"
        })
    });
}

interface getKeyListInterface {
    name: string,
    namespace: string
}
export function getKeyList({ name, namespace }: getKeyListInterface) {
    let headers = new Headers();
    addAuthentication(headers);
    return fetch(`joe2/osslist?name=${name}&namespace=${namespace}`, {
        method: "GET",
        headers
    })
}

interface searchImgInterface {
    name: string,
    namespace: string,
    tag: string
}
export function searchImg({ name, namespace, tag }: searchImgInterface) {
    let headers = new Headers();
    addAuthentication(headers);
    return fetch(`joe2/osslist?name=${name}&namespace=${namespace}&tag=${tag}`, {
        method: "GET",
        headers
    });
}

interface getObjectInterface {
    name: string,
    namespace: string,
    key: string
}
export function getObject({ name, namespace, key }: getObjectInterface) {
    let headers = new Headers();
    addAuthentication(headers);
    return fetch(`joe2/oss?name=${name}&namespace=${namespace}&key=${key}`, {
        method: "GET",
        headers
    });
}

interface updateObjectPvcInterface {
    name: string,
    namespace: string,
    maxobject: string,
    maxgbsize: string
}

function updateObjectPvc(props: updateObjectPvcInterface) {

}

export function checkShouldWait() {
    return fetch("joe/api/wait", {
        method: "GET"
    });
}