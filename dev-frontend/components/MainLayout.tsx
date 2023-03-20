import { ReactElement, useRef, useState } from "react";
import { Menu } from "antd";
import styles from "@/styles/components/MainLayout.module.scss"
import { useRouter } from "next/router";
import Link from "next/link";
import clsx from "clsx";

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

export default function MainLayout({children}: MainLayout) {
    let router = useRouter();
    
    let [isScrolling, setIsScrolling] = useState(false);
    
    let scrollClass: {[key: string]: boolean} = {};
    scrollClass[styles.scrolling] = isScrolling;
    
    let timeoutId = useRef<null | ReturnType<typeof setTimeout>>(null);
    let scrollDisappearTime = 1500;
    
    function handleScroll() {
        setIsScrolling(true);
        if (timeoutId.current) {
            clearTimeout(timeoutId.current);
        }
        timeoutId.current = setTimeout(() => {
            setIsScrolling(false);
        }, scrollDisappearTime);
    }
    return (
        <div className={styles.outerContainer}>
            <div className={styles.header}>
                <Link href={"/"}>
                    <p>
                        CubeUniverse
                    </p>
                </Link>
            </div>
            <div className={styles.bottomContainer}>
                <div 
                    className={clsx(
                        styles.sider,
                        scrollClass
                    )} 
                    onScroll={() => handleScroll()}
                >
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
                </div>
                <div 
                    className={clsx(
                        styles.contentContainer,
                        scrollClass
                    )} 
                >
                    {children}
                </div>
            </div>
        </div>
    )
}