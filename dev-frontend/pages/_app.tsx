import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import {
    Box,
    Button,
    Center,
    MantineProvider,
    Title,
    Paper,
    SimpleGrid,
    Stack,
    PasswordInput,
    Loader,
} from "@mantine/core"
import { NextPage } from 'next'
import { ReactElement, ReactNode, useEffect, useRef, useState } from 'react'
import { emotionCache } from '@/emotionCache'
import DataProvider from '@/components/DataProvider'
import { Notifications } from '@mantine/notifications'
import Head from 'next/head'
import { checkShouldWait, login } from '@/apis'
import { noti } from '@/utils/noti'
import { authentication } from '@/storage'
import Image from 'next/image'
import cubeUniverse from "../public/logo.png"
import { useForm } from '@mantine/form'
import { atom } from "signia";

export let wsUrl = atom("ws", "");

export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
    getLayout?: (page: ReactElement) => ReactNode
}

type AppPropsWithLayout = AppProps & {
    Component: NextPageWithLayout
}

export default function MyApp({ Component, pageProps }: AppPropsWithLayout) {
    // Use the layout defined at the page level, if available
    const getLayout = Component.getLayout ?? ((page) => page);
    let [isLogin, setIsLogin] = useState(false);
    let [shouldWait, setShouldWait] = useState(true);
    let intervalId = useRef<null | NodeJS.Timer>(null);
    useEffect(() => {
        intervalId.current = setInterval(() => {
            checkShouldWait()
                .then(e => e.status === 404 ? Promise.resolve() : Promise.reject()
                )
                .then(() => {
                    intervalId.current && clearInterval(intervalId.current)
                    setShouldWait(false);
                })
                .catch(() => { });
        }, 1000);
        return () => {
            intervalId.current && clearInterval(intervalId.current)
        };
    }, []);

    return (
        <>
            <Head>
                <title>CubeUniverse</title>
            </Head>
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
                {shouldWait ? <Wait />
                    : isLogin
                        ? (
                            <DataProvider>
                                {getLayout(<Component {...pageProps} />)}
                            </DataProvider>
                        )
                        : <Login setLogin={(x: boolean) => setIsLogin(x)} />}
            </MantineProvider>
        </>
    )
}


function Login({ setLogin }: { setLogin: (x: boolean) => void }) {
    let form = useForm({
        initialValues: {
            password: "",
        },
        validate: {
            password: val => val === "" ? "请输入密码" : null,
        }
    })
    return (
        <Box h={"100vh"} sx={{
            background: "linear-gradient(135deg, hsla(209, 79%, 81%, 1) 2%, hsla(266, 7%, 53%, 1) 99%)",
        }}>
            <Center h={"100%"}>
                <Paper
                    w={600}
                    h={300}
                    shadow='xl'
                    radius={"lg"}
                    sx={{
                        background: "rgba(255, 255, 255, 0.2)",
                        borderRadius: "16px",
                        boxShadow: "0 4px 30px rgba(0, 0, 0, 0.1)",
                        backdropFilter: "blur(5px)",
                    }}>
                    <SimpleGrid h={"100%"} cols={2}>
                        <Box h={"100%"} pl={"xl"}>
                            <Center h={"100%"}>
                                <Stack>
                                    <Center>
                                        <Image
                                            src={cubeUniverse}
                                            alt='Cube Universe'
                                            height={100}
                                            style={{ objectFit: "contain" }}
                                        />
                                    </Center>
                                    <Title
                                        order={1}
                                        variant="gradient"
                                        gradient={{ from: "#4fb9e3", to: "#032d81", deg: 30 }}
                                    >CubeUniverse</Title>
                                </Stack>
                            </Center>
                        </Box>
                        <Box
                            h={"100%"}>
                            <Stack h={"100%"} justify='center'>
                                <form
                                    onSubmit={form.onSubmit(vals => {
                                        login(vals.password)
                                            .then(e => {
                                                let url = e.url.replace("http", "ws");
                                                let arr = url.split("30400");
                                                url = [arr[0], "30401/api/storage/pvcws"].join();
                                                wsUrl.set(url);
                                                return e.json();
                                            })
                                            .then(e => e.code === 200 ? Promise.resolve(e) : Promise.reject(e))
                                            .then(e => {
                                                authentication.set(e.data.token);
                                                setLogin(true);
                                                noti(
                                                    'Success',
                                                    "Success",
                                                    "登陆成功");
                                            })
                                            .catch(e => {
                                                noti(
                                                    'Error',
                                                    "Error",
                                                    "密码错误");
                                            });
                                    })}>
                                    <Stack p={"xs"} h={"100%"} justify='center'>
                                        <PasswordInput
                                            {...form.getInputProps("password")}
                                            label={<Title order={6} italic>密码</Title>} />
                                        <Button
                                            type='submit'
                                            variant="gradient"
                                            gradient={{ from: 'teal', to: 'blue', deg: 60 }}>
                                            登录
                                        </Button>
                                    </Stack>
                                </form>
                            </Stack>
                        </Box>
                    </SimpleGrid>
                </Paper>
            </Center>
        </Box>
    )
}

function Wait() {
    return (
        <Box h={"100vh"} sx={{
            background: "linear-gradient(135deg, hsla(209, 79%, 81%, 1) 2%, hsla(266, 7%, 53%, 1) 99%)",
        }}>
            <Center h={"100%"}>
                <Stack>
                    <Title>等待集群建立……</Title>
                    <Center>
                        <Loader size={'xl'} />
                    </Center>
                </Stack>
            </Center>
        </Box>
    )
}