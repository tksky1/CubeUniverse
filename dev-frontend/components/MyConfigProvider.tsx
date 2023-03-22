import { ReactElement } from "react"
import { ConfigProvider } from "antd"

interface MyConfigProvider {
    children: ReactElement 
}

let myTheme = {
    token: {
      colorPrimary: "#0440a4",
    }
}

export default function MyConfigProvider({children}: MyConfigProvider) {
    return <ConfigProvider theme={myTheme}>{children}</ConfigProvider>
}