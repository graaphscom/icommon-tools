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
import {
  addressBook,
  addressCard,
  alignCenter,
} from '@icommon/fontawesome/solid';
import {
  __3dRotation24px,
  accessibleForward24px,
  accountBalanceSharp24px,
} from '@icommon/material/action';
import {
  addAlert24px,
  notificationImportant24px,
  warningSharp24px,
} from '@icommon/material/alert';
import { arrowBoth24, blocked24, briefcase24 } from '@icommon/octicons';
import { activityLog, borderDashed, button } from '@icommon/radixui';
import {
  ancientGateFill,
  bankFill,
  homeFill,
} from '@icommon/remixicon/buildings';
import {
  anticlockwiseFill,
  brush2Fill,
  cropFill,
} from '@icommon/remixicon/design';
import { __0Plus, __12Plus, __500px } from '@icommon/unicons/line';
import { angleDoubleDown, apps, bookmark } from '@icommon/unicons/monochrome';

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
          <Icommon node={addressBook} width="32" height="32" />
          <Icommon node={addressCard} width="32" height="32" />
          <Icommon node={alignCenter} width="32" height="32" />
        </section>
      </section>
      <section>
        <h2>material</h2>
        <section>
          <h3>action</h3>
          <Icommon node={__3dRotation24px} />
          <Icommon node={accessibleForward24px} />
          <Icommon node={accountBalanceSharp24px} />
        </section>
        <section>
          <h3>alert</h3>
          <Icommon node={addAlert24px} />
          <Icommon node={notificationImportant24px} />
          <Icommon node={warningSharp24px} />
        </section>
      </section>
      <section>
        <h2>octicons</h2>
        <Icommon node={arrowBoth24} />
        <Icommon node={blocked24} />
        <Icommon node={briefcase24} />
      </section>
      <section>
        <h2>radixui</h2>
        <Icommon node={activityLog} />
        <Icommon node={borderDashed} />
        <Icommon node={button} />
      </section>
      <section>
        <h2>remixicon</h2>
        <section>
          <h3>buildings</h3>
          <Icommon node={ancientGateFill} width="32" height="32" />
          <Icommon node={bankFill} width="32" height="32" />
          <Icommon node={homeFill} width="32" height="32" />
        </section>
        <section>
          <h3>design</h3>
          <Icommon node={anticlockwiseFill} width="32" height="32" />
          <Icommon node={brush2Fill} width="32" height="32" />
          <Icommon node={cropFill} width="32" height="32" />
        </section>
      </section>
      <section>
        <h2>unicons</h2>
        <section>
          <h3>line</h3>
          <Icommon node={__0Plus} width="32" height="32" />
          <Icommon node={__12Plus} width="32" height="32" />
          <Icommon node={__500px} width="32" height="32" />
        </section>
        <section>
          <h3>monochrome</h3>
          <Icommon node={angleDoubleDown} width="32" height="32" />
          <Icommon node={apps} width="32" height="32" />
          <Icommon node={bookmark} width="32" height="32" />
        </section>
      </section>
    </section>
  );
}
