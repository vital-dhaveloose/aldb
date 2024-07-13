import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component';
import { provideFluentDesignSystem, fluentCard, fluentButton, fluentTextField, fluentTreeItem, fluentTreeView } from '@fluentui/web-components';


bootstrapApplication(AppComponent, appConfig)
.catch((err) => console.error(err));

provideFluentDesignSystem().register(
  fluentCard(),
  fluentButton(),
  fluentTextField(),
  fluentTreeView(),
  fluentTreeItem()
);