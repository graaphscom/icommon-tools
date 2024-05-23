import { Icommon } from '@icommon/components/icommon';
import { bxl500px, bxl99designs, bxlAdobe } from '@icommon/boxicons/logos';
import { bxAbacus, bxAccessibility, bxAlarm } from '@icommon/boxicons/regular';
import {
  bxsAddToQueue,
  bxsAdjust,
  bxsAdjustAlt,
} from '@icommon/boxicons/solid';
import { activity, ban, bag } from '@icommon/bytesize';
import { key32Filled } from '@icommon/fluentui/key';
import { airplaneLanding24Filled } from '@icommon/fluentui/airplaneLanding';
import { addSquareMultiple24Filled } from '@icommon/fluentui/addSquareMultiple';
import { node, adn, aws } from '@icommon/fontawesome/brands';
import { map, bell, eye } from '@icommon/fontawesome/regular';

export default function Home() {
  return (
    <section>
      <h1>icommon</h1>
      <section>
        <h2>boxicons</h2>
        <section>
          <h3>logos</h3>
          <Icommon node={bxl99designs} />
          <Icommon node={bxl500px} />
          <Icommon node={bxlAdobe} />
        </section>
        <section>
          <h3>regular</h3>
          <Icommon node={bxAbacus} />
          <Icommon node={bxAccessibility} />
          <Icommon node={bxAlarm} />
        </section>
        <section>
          <h3>solid</h3>
          <Icommon node={bxsAddToQueue} />
          <Icommon node={bxsAdjust} />
          <Icommon node={bxsAdjustAlt} />
        </section>
      </section>
      <section>
        <h2>bytesize</h2>
        <Icommon node={activity} />
        <Icommon node={bag} />
        <Icommon node={ban} />
      </section>
      <section>
        <h2>fluentui</h2>
        <Icommon node={key32Filled} />
        <Icommon node={airplaneLanding24Filled} />
        <Icommon node={addSquareMultiple24Filled} />
      </section>
      <section>
        <h2>fontawesome</h2>
        <section>
          <h3>brands</h3>
          <Icommon node={node} width="32" height="32" />
          <Icommon node={adn} width="32" height="32" />
          <Icommon node={aws} width="32" height="32" />
        </section>
        <section>
          <h3>regular</h3>
          <Icommon node={map} width="32" height="32" />
          <Icommon node={bell} width="32" height="32" />
          <Icommon node={eye} width="32" height="32" />
        </section>
        <section>
          <h3>solid</h3>
        </section>
      </section>
      <h2>material</h2>
      <h2>octicons</h2>
      <h2>radixui</h2>
      <h2>remixicon</h2>
      <h2>unicons</h2>
    </section>
  );
}
