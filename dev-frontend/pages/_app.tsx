import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { MantineProvider } from "@mantine/core"
import { NextPage } from 'next'
import { ReactElement, ReactNode } from 'react'
import { emotionCache } from '@/emotionCache'
import DataProvider from '@/components/DataProvider'
import { data, DataContext } from '@/components/DataContext'
import { Notifications } from '@mantine/notifications'

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
    <MantineProvider
      withGlobalStyles
      withNormalizeCSS
      emotionCache={emotionCache}
      theme={{
        colorScheme: "light",
        colors: {
        }
      }}>
      <Notifications />
      <DataProvider>
        {getLayout(<Component {...pageProps} />)}
      </DataProvider>
    </MantineProvider>
  )
}