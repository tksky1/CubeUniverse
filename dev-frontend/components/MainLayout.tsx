import { ReactElement, useRef, useState } from "react";
import MyConfigProvider from "./MyConfigProvider";
import { Menu } from "antd";
import { useRouter } from "next/router";
import Link from "next/link";
import { Title, createStyles, ScrollArea, Flex, Container, Text, useMantineTheme, Image, Group, Button } from "@mantine/core"
import { logo } from "@/utils/logostr";
import { token } from "@/utils/token";
import { login } from "@/apis";

interface MainLayout {
    children: ReactElement
}

let items = [
    {
        title: "仪表盘",
        key: "/",
        label: "仪表盘"
    },
    {
        title: "存储池",
        key: "/cephPool",
        label: "存储池"
    },
    {
        title: "节点",
        key: "/node",
        label: "节点"
    },
    {
        title: "日志",
        key: "/log",
        label: "日志"
    },
    {
        type: "group" as const,
        label: " 存储 ",
        children: [
            {
                title: "文件存储",
                key: "/fileStorage",
                label: "文件存储"
            },
            {
                title: "块存储",
                key: "/blockStorage",
                label: "块存储"
            },
            {
                title: "对象存储",
                key: "/objectStorage",
                label: "对象存储"
            },
        ]
    }
]

const useStyles = createStyles(theme => {
    return {
        outerContainer: {
            height: "100vh",
        },
        bottomContainer: {
            height: "calc(100vh - 50px)",
            display: "flex"
        },
        scrollAreaViewport: {
            ' > div': {
                height: "100%"
            }
        }
    }
})

export default function MainLayout({ children }: MainLayout) {
    let { classes } = useStyles();
    let router = useRouter();
    let theme = useMantineTheme();
    return (
        (
            <div
                className={
                    classes.outerContainer
                }
            >
                <Flex
                    sx={theme => ({
                        height: 50,
                        boxShadow: `inset 0 -7px 7px -5px ${theme.colors.blue[4]}`,
                        backgroundColor: "#000000"
                    })}
                    justify="flex-start"
                    align="center"
                    pl={10}
                >
                    <Link href={"/"}>
                        <Group spacing={5}>
                            <Image
                                height={30}
                                width={30}
                                fit="contain"
                                src={logo()} />
                            <Title
                                order={1}
                                variant="gradient"
                                gradient={{ from: "#4fb9e3", to: "#032d81", deg: 30 }}
                            >CubeUniverse</Title>
                        </Group>
                    </Link>
                </Flex>
                <div className={classes.bottomContainer}>
                    <ScrollArea
                        sx={theme => ({
                            height: "100%",
                            minWidth: "200px",
                        })}
                        classNames={{
                            viewport: classes.scrollAreaViewport
                        }}
                        type="scroll"
                    >
                        <MyConfigProvider>
                            <Menu
                                theme="light"
                                selectedKeys={[router.asPath]}
                                items={items}
                                mode="inline"
                                defaultOpenKeys={["/objectStorageGroup"]}
                                onClick={(clickedItem) => router.push(clickedItem.key)}
                                style={{
                                    height: "100%",
                                    width: "100%",
                                    borderStyle: "none",
                                }}
                            />
                        </MyConfigProvider>
                    </ScrollArea>
                    <ScrollArea
                        sx={_theme => ({
                            flexGrow: 1,
                            height: "100%",
                        })}
                        classNames={{
                            viewport: classes.scrollAreaViewport
                        }}
                        type="scroll"
                    >
                        {children}
                    </ScrollArea>
                </div>
            </div>
        )
    )
}