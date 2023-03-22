import { ReactElement, useRef, useState } from "react";
import MyConfigProvider from "./MyConfigProvider";
import { Menu } from "antd";
import { useRouter } from "next/router";
import Link from "next/link";
import { Title, createStyles, ScrollArea, Flex, Container } from "@mantine/core"

interface MainLayout {
    children: ReactElement
}

let items = [
    {
        title: "page1",
        key: "/page1",
        label: "page1"
    },
    {
        title: "page2",
        key: "/page2",
        label: "page2"
    },
    {
        children: [
            {
                title: "page3",
                key: "/nest/page3",
                label: "page3"
            },
            {
                title: "page4",
                key: "/nest/page4",
                label: "page4"
            },
        ],
        key: "/nest",
        label: "nest"
    },
    {
        type: "group" as const,
        label: "group",
        children: [
            {
                title: "page5",
                key: "/group/page5",
                label: "page5"
            },
            {
                title: "page6",
                key: "/group/page6",
                label: "page6"
            },
        ]
    }
]

const useStyles = createStyles(_theme => {
    return {
        outerContainer: {
            height: "100vh"
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

export default function MainLayout({children}: MainLayout) {
    let { classes } = useStyles();
    let router = useRouter();
    
    return (
        <div className={classes.outerContainer}>
            <Flex
                sx={_theme => ({
                    height: 50,
                    boxShadow: 'inset 0 -5px 7px -5px #0440a4',
                })}
                justify="flex-start"
                align="center"
                pl={10}
            >
                <Link href={"/"}>
                    <Title 
                        order={1}
                        variant="gradient"
                        gradient={{from: "#4fb9e3", to: "#032d81", deg: 30}}
                    >CubeUniverse</Title>
                </Link>
            </Flex>
            <div className={classes.bottomContainer}>
                <ScrollArea 
                    sx={_theme => ({
                        height: "100%",
                        width: 200
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
                            onClick={(clickedItem) => router.push(clickedItem.key)}
                            style={{
                                height: "100%",
                                width: "100%",
                            }}
                        />
                    </MyConfigProvider>
                </ScrollArea>
                <ScrollArea
                    sx={_theme => ({
                        flexGrow: 1,
                        height: "100%",
                        backgroundColor: "#eeeeee",
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
}