import { ReactElement } from "react"
import { ConfigProvider } from "antd"
import { useMantineTheme } from "@mantine/core"

interface MyConfigProvider {
    children: ReactElement
}

let myTheme = {
    token: {
        colorPrimary: "#0440a4",
    }
}

export default function MyConfigProvider({ children }: MyConfigProvider) {
    let theme = useMantineTheme();
    return <ConfigProvider theme={myTheme}>{children}</ConfigProvider>
}