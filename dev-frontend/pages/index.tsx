import MainLayout, { Entry } from "@/components/MainLayout";
import { commonLayout } from "@/utils/commonLayout";
import { ReactElement } from "react";

export async function getStaticProps() {
  return {
    props: {
    }
  }
}

export default function Home() {
  return (
    <>
    <div className="flex items-center justify-center h-full">
      joejoe
    </div>
    </>
  )
}

Home.getLayout = commonLayout;