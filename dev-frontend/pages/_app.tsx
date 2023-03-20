import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { ConfigProvider, theme } from 'antd'
import { NextPage } from 'next'
import { ReactElement, ReactNode } from 'react'

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getLayout?: (page: ReactElement) => ReactNode
}

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout
}

let myTheme = {
  token: {
    colorPrimary: "#0440a4",
  }
}

export default function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  // Use the layout defined at the page level, if available
  const getLayout = Component.getLayout ?? ((page) => page)

  return (
    <ConfigProvider theme={myTheme}>
      {getLayout(<Component {...pageProps} />)}
    </ConfigProvider>
  )
}