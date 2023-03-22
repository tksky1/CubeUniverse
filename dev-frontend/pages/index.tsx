import { commonLayout } from "@/utils/commonLayout";
import { Card, Text, Metric, Flex, ProgressBar } from "@tremor/react";

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
        <Card className="max-w-xs mx-auto">
          <Text>Sales</Text>
          <Metric>$ 71,465</Metric>
          <Flex className="mt-4">
            <Text>32% of annual target</Text>
            <Text>$ 225,000</Text>
          </Flex>
          <ProgressBar percentageValue={32} className="mt-2" />
        </Card>
      </div>
    </>
  )
}

Home.getLayout = commonLayout;